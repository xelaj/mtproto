// @constructor Defines a [min](https://core.telegram.org/api/min) user that was seen in a certain message of a certain chat.
// @param peer The chat where the user was seen
// @param msg_id The message ID
// @param user_id The identifier of the user that was seen
inputPeerUserFromMessage#17bae2e6 peer:InputPeer msg_id:int user_id:int = InputPeer;

// @type Object describes the file type.
// @enum Unknown type.
storage.fileUnknown#aa963b05 = storage.FileType;
// @enum Part of a bigger file.
storage.filePartial#40bc6f52 = storage.FileType;
// @enum JPEG image. MIME type: image/jpeg.
storage.fileJpeg#7efe0e = storage.FileType;
// @enum GIF image. MIME type: image/gif.
storage.fileGif#cae1aadf = storage.FileType;
// @enum PNG image. MIME type: image/png.
storage.filePng#a4f63c0 = storage.FileType;
// @enum PDF document image. MIME type: application/pdf.
storage.filePdf#ae1e508d = storage.FileType;
// @enum Mp3 audio. MIME type: audio/mpeg.
storage.fileMp3#528a0677 = storage.FileType;
// @enum Quicktime video. MIME type: video/quicktime.
storage.fileMov#4b09ebbc = storage.FileType;
// @enum MPEG-4 video. MIME type: video/mp4.
storage.fileMp4#b3cea0e4 = storage.FileType;
// @enum WEBP image. MIME type: image/webp.
storage.fileWebp#1081464c = storage.FileType;