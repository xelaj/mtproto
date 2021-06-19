package main

// original could be found at https://github.com/umputun/feed-master/pull/37

import (
	"hash/maphash"
	"io"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/xelaj/go-dry"
	"github.com/xelaj/mtproto/telegram"

	utils "github.com/xelaj/mtproto/examples/example_utils"
)

const kb = 1024
const mb = kb * 1024

// Item for rss
type Item struct {
	Title     string    `xml:"title"`
	Link      string    `xml:"link"`
	Enclosure Enclosure `xml:"enclosure"`
}

// Enclosure element from item
type Enclosure struct {
	URL    string `xml:"url,attr"`
	Length int    `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

func main() {
	println("firstly, you need to authorize. after example 'auth', you will sign in")

	// helper variables
	appStorage := utils.PrepareAppStorageForExamples()
	sessionFile := filepath.Join(appStorage, "session.json")
	publicKeys := filepath.Join(appStorage, "tg_public_keys.pem")

	client, err := telegram.NewClient(telegram.ClientConfig{
		// where to store session configuration. must be set
		SessionFile: sessionFile,
		// host address of mtproto server. Actually, it can be any mtproxy, not only official
		ServerHost: "149.154.167.50:443",
		// public keys file is path to file with public keys, which you must get from https://my.telelgram.org
		PublicKeysFile:  publicKeys,
		AppID:           94575,                              // app id, could be find at https://my.telegram.org
		AppHash:         "a3406de8d171bb422bb6ddf3bbd800e2", // app hash, could be find at https://my.telegram.org
		InitWarnChannel: true,                               // if we want to get errors, otherwise, client.Warnings will be set nil
	})
	utils.ReadWarningsToStdErr(client.Warnings)
	dry.PanicIfErr(err)

	// authorize the bot
	_, err = client.AuthImportBotAuthorization(0, 94575, "a3406de8d171bb422bb6ddf3bbd800e2", "here_goes_the_token")
	dry.PanicIfErr(errors.Wrapf(err, "error authorizing with telegram bot"))

	// get reference for the channel we'll be posting to
	chanRef, err := getChannelReference(client, "here_goes_your_public_channel_name")
	dry.PanicIfErr(errors.Wrapf(err, "error retrieving channel metadata"))

	// here you should prepare an object to send
	item := Item{
		Title: "This is a test message to send",
		Link:  "https://example.org",
		Enclosure: Enclosure{
			URL:  "https://file-examples-com.github.io/uploads/2017/11/file_example_MP3_700KB.mp3",
			Type: "audio/mpeg3",
		},
	}

	// upload the file and send the message with the media
	err = sendMessageWithFile(client, item, chanRef, item.Title)
	dry.PanicIfErr(err)
}

// getChannelReference returns telegram channel metadata reference which
// is enough to send messages to that channel using the telegram API
func getChannelReference(client *telegram.Client, channelID string) (telegram.InputPeer, error) {
	channel, err := client.ContactsResolveUsername(channelID)
	if err != nil {
		return nil, err
	}
	return &telegram.InputPeerChannel{
		ChannelID:  channel.Chats[0].(*telegram.Channel).ID,
		AccessHash: channel.Chats[0].(*telegram.Channel).AccessHash,
	}, nil
}

// getContentLength uses HEAD request and called as a fallback in case of item.Enclosure.Length not populated
func getContentLength(url string) (int, error) {
	resp, err := http.Head(url) // nolint:gosec // URL considered safe
	if err != nil {
		return 0, errors.Wrapf(err, "can't HEAD %s", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.Errorf("non-200 status, %d", resp.StatusCode)
	}

	return int(resp.ContentLength), err
}

// downloadAudio returns resp.Body for provided URL
func downloadAudio(url string) (io.ReadCloser, error) {
	clientHTTP := &http.Client{Timeout: time.Minute}

	resp, err := clientHTTP.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// getInt64Hash generates int64 hash from the provided string, returns 0 in case of error
func getInt64Hash(s string) int64 {
	hash := maphash.Hash{}
	_, _ = hash.Write([]byte(s))
	return int64(hash.Sum64())
}

// getFilenameByURL returns filename from a given URL
func getFilenameByURL(url string) string {
	_, filename := path.Split(url)
	return filename
}

// uploadFileToTelegram uploads file to telegram API returns number of file parts it uploaded
func uploadFileToTelegram(client *telegram.Client, r io.Reader, fileID int64, fileLength int) (int32, error) {
	var fileParts []int32
	// 512kb is magic number from https://core.telegram.org/api/files, you can't set bigger chunks
	chunkSize := 512 * kb
	buf := make([]byte, chunkSize)
	approximateChunks := int32(fileLength/chunkSize + 1)
	var err error
	var copyBytes int
	for err != io.EOF && err != io.ErrUnexpectedEOF {
		copyBytes, err = io.ReadFull(r, buf)
		if err != io.EOF && err != io.ErrUnexpectedEOF && err != nil {
			return 0, errors.Wrapf(err, "error reading the file chunk for upload")
		}
		// don't send zero-filled buffer part in case that's the last chunk of file
		if err == io.ErrUnexpectedEOF {
			buf = buf[:copyBytes]
		}

		filePartID := int32(len(fileParts))
		_, uploadErr := client.UploadSaveBigFilePart(fileID, filePartID, approximateChunks, buf)
		if uploadErr != nil {
			return 0, errors.Wrapf(uploadErr, "error uploading the file using telegram API")
		}
		fileParts = append(fileParts, filePartID)
	}
	return int32(len(fileParts)), nil
}

func sendMessageWithFile(client *telegram.Client, item Item, chanRef telegram.InputPeer, msg string) error {
	contentLength, err := getContentLength(item.Enclosure.URL)
	if err != nil {
		return errors.Wrapf(err, "can't get length for %s", item.Enclosure.URL)
	}

	log.Printf("[DEBUG] start uploading audio %s (%dMb)", item.Enclosure.URL, contentLength/mb)
	httpBody, err := downloadAudio(item.Enclosure.URL)
	if err != nil {
		return errors.Wrapf(err, "error retrieving audio")
	}
	defer httpBody.Close()

	fileID := getInt64Hash(item.Enclosure.URL)
	fileChunks, err := uploadFileToTelegram(client, httpBody, fileID, contentLength)
	if err != nil {
		return errors.Wrapf(err, "error uploading the file")
	}

	mimeType := item.Enclosure.Type
	if mimeType == "" {
		mimeType = "audio/mpeg"
	}

	_, err = client.MessagesSendMedia(&telegram.MessagesSendMediaParams{
		Peer: chanRef,
		Media: &telegram.InputMediaUploadedDocument{
			MimeType: mimeType,
			Attributes: []telegram.DocumentAttribute{
				&telegram.DocumentAttributeAudio{Title: item.Title},
				&telegram.DocumentAttributeFilename{FileName: getFilenameByURL(item.Enclosure.URL)},
			},
			File: &telegram.InputFileBig{
				ID:    fileID,
				Parts: fileChunks,
				Name:  getFilenameByURL(item.Enclosure.URL),
			},
		},
		RandomID: fileID,
		Message:  msg,
	})
	if err != nil {
		return errors.Wrapf(err, "error uploading message to channel")
	}

	return nil
}
