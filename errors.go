// also why gocritic detects false positive, but if i write explanation, golangci-lint throws error that description expected as lintrer??? //TODO
//nolint: lll
package mtproto

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/xelaj/mtproto/internal/mtproto/objects"
)

type ErrResponseCode struct {
	Code           int
	Message        string
	Description    string
	AdditionalInfo any // some errors has additional data like timeout seconds, dc id etc.
}

func RpcErrorToNative(r *objects.RpcError) error {
	nativeErrorName, additionalData := TryExpandError(r.ErrorMessage)

	desc, ok := errorMessages[nativeErrorName]
	if !ok {
		desc = nativeErrorName
	}

	if additionalData != nil {
		desc = fmt.Sprintf(desc, additionalData)
	}

	return &ErrResponseCode{
		Code:           int(r.ErrorCode),
		Message:        nativeErrorName,
		Description:    desc,
		AdditionalInfo: additionalData,
	}
}

type prefixSuffix struct {
	prefix string
	suffix string
	kind   reflect.Kind // int string bool etc.
}

var specificErrors = []prefixSuffix{
	{"EMAIL_UNCONFIRMED_", "", reflect.Int},
	{"FILE_MIGRATE_", "", reflect.Int},
	{"FILE_PART_", "_MISSING", reflect.Int},
	{"FLOOD_TEST_PHONE_WAIT_", "", reflect.Int},
	{"FLOOD_WAIT_", "", reflect.Int},
	{"INTERDC_", "_CALL_ERROR", reflect.Int},
	{"INTERDC_", "_CALL_RICH_ERROR", reflect.Int},
	{"NETWORK_MIGRATE_", "", reflect.Int},
	{"PASSWORD_TOO_FRESH_", "", reflect.Int},
	{"PHONE_MIGRATE_", "", reflect.Int},
	{"SESSION_TOO_FRESH_", "", reflect.Int},
	{"SLOWMODE_WAIT_", "", reflect.Int},
	{"STATS_MIGRATE_", "", reflect.Int},
	{"TAKEOUT_INIT_DELAY_", "", reflect.Int},
	{"USER_MIGRATE_", "", reflect.Int},
}

func TryExpandError(errStr string) (nativeErrorName string, additionalData any) {
	var choosedPrefixSuffix *prefixSuffix

	for _, errCase := range specificErrors {
		if strings.HasPrefix(errStr, errCase.prefix) && strings.HasSuffix(errStr, errCase.suffix) {
			choosedPrefixSuffix = &errCase //nolint:gosec cause we need nil if not found
			break
		}
	}

	if choosedPrefixSuffix == nil {
		return errStr, nil // common error, returning
	}

	nativeErrorName = choosedPrefixSuffix.prefix + "X" + choosedPrefixSuffix.suffix
	trimmedData := strings.TrimSuffix(strings.TrimPrefix(errStr, choosedPrefixSuffix.prefix), choosedPrefixSuffix.suffix)

	switch v := choosedPrefixSuffix.kind; v { //nolint:exhaustive others will panic
	case reflect.Int:
		var err error
		additionalData, err = strconv.Atoi(trimmedData)
		check(errors.Wrap(err, "error of parsing expected int value"))

	case reflect.String:
		additionalData = trimmedData

	default:
		panic("couldn't parse this type: " + v.String())
	}

	return nativeErrorName, additionalData
}

func (e *ErrResponseCode) Error() string {
	return fmt.Sprintf("%s (code %d)", e.Description, e.Code)
}

// gathered all errors from all methods. don't have reference in docs at all
var errorMessages = map[string]string{
	"ABOUT_TOO_LONG":                      "About string too long",
	"ACCESS_TOKEN_EXPIRED":                "Bot token expired",
	"ACCESS_TOKEN_INVALID":                "The provided token is not valid",
	"ACTIVE_USER_REQUIRED":                "The method is only available to already activated users",
	"ADMINS_TOO_MUCH":                     "Too many admins",
	"ADMIN_RANK_EMOJI_NOT_ALLOWED":        "Emoji are not allowed in admin titles or ranks",
	"ADMIN_RANK_INVALID":                  "The given admin title or rank was invalid (possibly larger than 16 characters)",
	"API_ID_INVALID":                      "API ID invalid",
	"API_ID_PUBLISHED_FLOOD":              "This API id was published somewhere, you can't use it now",
	"ARTICLE_TITLE_EMPTY":                 "The title of the article is empty",
	"AUTH_BYTES_INVALID":                  "The provided authorization is invalid",
	"AUTH_KEY_DUPLICATED":                 "The authorization key (session file) was used under two different IP addresses simultaneously, and can no longer be used. Use the same session exclusively, or use different sessions",
	"AUTH_KEY_INVALID":                    "Auth key invalid",
	"AUTH_KEY_PERM_EMPTY":                 "The method is unavailable for temporary authorization key, not bound to permanent",
	"AUTH_KEY_UNREGISTERED":               "The key is not registered in the system",
	"AUTH_RESTART":                        "Restart the authorization process",
	"AUTH_TOKEN_ALREADY_ACCEPTED":         "The authorization token was already used",
	"AUTH_TOKEN_EXPIRED":                  "The authorization token has expired",
	"AUTH_TOKEN_INVALID":                  "An invalid authorization token was provided",
	"AUTH_TOKEN_INVALIDX":                 "The specified auth token is invalid",
	"BANNED_RIGHTS_INVALID":               "You provided some invalid flags in the banned rights",
	"BOT_CHANNELS_NA":                     "Bots can't edit admin privileges",
	"BOT_COMMAND_DESCRIPTION_INVALID":     "The command description was empty, too long or had invalid characters used",
	"BOT_DOMAIN_INVALID":                  "Bot domain invalid",
	"BOT_GROUPS_BLOCKED":                  "This bot can't be added to groups",
	"BOT_INLINE_DISABLED":                 "This bot can't be used in inline mode",
	"BOT_INVALID":                         "This is not a valid bot",
	"BOT_METHOD_INVALID":                  "The API access for bot users is restricted. The method you tried to invoke cannot be executed as a bot",
	"BOT_MISSING":                         "This method can only be run by a bot",
	"BOT_PAYMENTS_DISABLED":               "This method can only be run by a bot",
	"BOT_POLLS_DISABLED":                  "You cannot create polls under a bot account",
	"BOT_RESPONSE_TIMEOUT":                "The bot did not answer to the callback query in time",
	"BOTS_TOO_MUCH":                       "There are too many bots in this chat/channel",
	"BROADCAST_FORBIDDEN":                 "The request cannot be used in broadcast channels",
	"BROADCAST_ID_INVALID":                "The channel is invalid",
	"BROADCAST_PUBLIC_VOTERS_FORBIDDEN":   "You can't forward polls with public voters",
	"BROADCAST_REQUIRED":                  "The request can only be used with a broadcast channel",
	"BUTTON_DATA_INVALID":                 "The data of one or more of the buttons you provided is invalid",
	"BUTTON_TYPE_INVALID":                 "The type of one of the buttons you provided is invalid",
	"BUTTON_URL_INVALID":                  "Button URL invalid",
	"CALL_ALREADY_ACCEPTED":               "The call was already accepted",
	"CALL_ALREADY_DECLINED":               "The call was already declined",
	"CALL_OCCUPY_FAILED":                  "The call failed because the user is already making another call",
	"CALL_PEER_INVALID":                   "The provided call peer object is invalid",
	"CALL_PROTOCOL_FLAGS_INVALID":         "Call protocol flags invalid",
	"CDN_METHOD_INVALID":                  "You can't call this method in a CDN DC",
	"CHANNEL_INVALID":                     "The provided channel is invalid",
	"CHANNEL_PRIVATE":                     "The channel specified is private and you lack permission to access it. Another reason may be that you were banned from it",
	"CHANNEL_PUBLIC_GROUP_NA":             "channel/supergroup not available",
	"CHANNEL_TOO_LARGE":                   "Channel is too large to be deleted; this error is issued when trying to delete channels with more than 1000 members (subject to change)",
	"CHANNELS_ADMIN_LOCATED_TOO_MUCH":     "Returned if both the check_limit and the by_location flags are set and the user has reached the limit of public geogroups",
	"CHANNELS_ADMIN_PUBLIC_TOO_MUCH":      "You're admin of too many public channels, make some channels private to change the username of this channel",
	"CHANNELS_TOO_MUCH":                   "You have joined too many channels/supergroups",
	"CHAT_ABOUT_NOT_MODIFIED":             "About text has not changed",
	"CHAT_ABOUT_TOO_LONG":                 "Chat about too long",
	"CHAT_ADMIN_INVITE_REQUIRED":          "You do not have the rights to do this",
	"CHAT_ADMIN_REQUIRED":                 "You must be an admin in this chat to do this",
	"CHAT_FORBIDDEN":                      "You cannot write in this chat",
	"CHAT_ID_EMPTY":                       "The provided chat ID is empty",
	"CHAT_ID_INVALID":                     "The provided chat id is invalid",
	"CHAT_INVALID":                        "Invalid chat",
	"CHAT_LINK_EXISTS":                    "The chat is linked to a channel and cannot be used in that request",
	"CHAT_NOT_MODIFIED":                   "The pinned message wasn't modified",
	"CHAT_RESTRICTED":                     "You can't send messages in this chat, you were restricted",
	"CHAT_SEND_GIFS_FORBIDDEN":            "You can't send gifs in this chat",
	"CHAT_SEND_INLINE_FORBIDDEN":          "You cannot send inline results in this chat",
	"CHAT_SEND_MEDIA_FORBIDDEN":           "You can't send media in this chat",
	"CHAT_SEND_POLL_FORBIDDEN":            "You can't send polls in this chat",
	"CHAT_SEND_STICKERS_FORBIDDEN":        "You can't send stickers in this chat",
	"CHAT_TITLE_EMPTY":                    "No chat title provided",
	"CHAT_WRITE_FORBIDDEN":                "You can't write in this chat",
	"CODE_EMPTY":                          "The provided code is empty",
	"CODE_HASH_INVALID":                   "Code hash invalid",
	"CODE_INVALID":                        "Code invalid",
	"CONNECTION_API_ID_INVALID":           "The provided API id is invalid",
	"CONNECTION_APP_VERSION_EMPTY":        "App version is empty",
	"CONNECTION_DEVICE_MODEL_EMPTY":       "Device model empty",
	"CONNECTION_LANG_PACK_INVALID":        "Language pack invalid",
	"CONNECTION_LAYER_INVALID":            "The very first request must always be InvokeWithLayerRequest",
	"CONNECTION_NOT_INITED":               "Connection not initialized",
	"CONNECTION_SYSTEM_EMPTY":             "Connection system empty",
	"CONNECTION_SYSTEM_LANG_CODE_EMPTY":   "The system language string was empty during connection",
	"CONTACT_ADD_MISSING":                 "Contact to add is missing",
	"CONTACT_ID_INVALID":                  "The provided contact ID is invalid",
	"CONTACT_NAME_EMPTY":                  "The provided contact name cannot be empty",
	"CONTACT_REQ_MISSING":                 "Missing contact request",
	"DATA_INVALID":                        "Encrypted data invalid",
	"DATA_JSON_INVALID":                   "The provided JSON data is invalid",
	"DATA_TOO_LONG":                       "Data too long",
	"DATE_EMPTY":                          "Date empty",
	"DC_ID_INVALID":                       "The provided DC ID is invalid",
	"DH_G_A_INVALID":                      "g_a invalid",
	"EMAIL_HASH_EXPIRED":                  "The email hash expired and cannot be used to verify it",
	"EMAIL_INVALID":                       "The given email is invalid",
	"EMAIL_UNCONFIRMED":                   "Email unconfirmed",
	"EMAIL_VERIFY_EXPIRED":                "The verification email has expired",
	"EMOTICON_EMPTY":                      "The emoticon field cannot be empty",
	"EMOTICON_INVALID":                    "The specified emoticon cannot be used or was not a emoticon",
	"ENCRYPTED_MESSAGE_INVALID":           "Encrypted message invalid",
	"ENCRYPTION_ALREADY_ACCEPTED":         "Secret chat already accepted",
	"ENCRYPTION_ALREADY_DECLINED":         "The secret chat was already declined",
	"ENCRYPTION_DECLINED":                 "The secret chat was declined",
	"ENCRYPTION_ID_INVALID":               "The provided secret chat ID is invalid",
	"ENCRYPTION_OCCUPY_FAILED":            "TDLib developer claimed it is not an error while accepting secret chats and 500 is used instead of 420",
	"ENTITIES_TOO_LONG":                   "It is no longer possible to send such long data inside entity tags (for example inline text URLs)",
	"ENTITY_MENTION_USER_INVALID":         "You mentioned an invalid user",
	"ERROR_TEXT_EMPTY":                    "The provided error message is empty",
	"EXPORT_CARD_INVALID":                 "Provided card is invalid",
	"EXTERNAL_URL_INVALID":                "External URL invalid",
	"FIELD_NAME_EMPTY":                    "The field with the name FIELD_NAME is missing",
	"FIELD_NAME_INVALID":                  "The field with the name FIELD_NAME is invalid",
	"FILE_ID_INVALID":                     "The provided file id is invalid",
	"FILE_PART_0_MISSING":                 "File part 0 missing",
	"FILE_PART_EMPTY":                     "The provided file part is empty",
	"FILE_PART_INVALID":                   "The file part number is invalid",
	"FILE_PART_LENGTH_INVALID":            "The length of a file part is invalid",
	"FILE_PART_SIZE_CHANGED":              "The file part size (chunk size) cannot change during upload",
	"FILE_PART_SIZE_INVALID":              "The provided file part size is invalid",
	"FILE_PART_TOO_BIG":                   "The uploaded file part is too big",
	"FILE_PARTS_INVALID":                  "The number of file parts is invalid",
	"FILE_REFERENCE_EMPTY":                "The file reference must exist to access the media and it cannot be empty",
	"FILE_REFERENCE_EXPIRED":              "File reference expired, it must be refetched as described in https://core.telegram.org/api/file_reference",
	"FILEREF_UPGRADE_NEEDED":              "The file reference needs to be refreshed before being used again",
	"FILTER_ID_INVALID":                   "The specified filter ID is invalid",
	"FIRSTNAME_INVALID":                   "The first name is invalid",
	"FOLDER_ID_EMPTY":                     "The folder you tried to delete was already empty",
	"FOLDER_ID_INVALID":                   "The folder you tried to use was not valid",
	"FRESH_CHANGE_ADMINS_FORBIDDEN":       "You were just elected admin, you can't add or modify other admins yet",
	"FRESH_CHANGE_PHONE_FORBIDDEN":        "Recently logged-in users cannot use this request",
	"FRESH_RESET_AUTHORISATION_FORBIDDEN": "You can't logout other sessions if less than 24 hours have passed since you logged on the current session",
	"FROM_MESSAGE_BOT_DISABLED":           "Bots can't use fromMessage min constructors",
	"GAME_BOT_INVALID":                    "You cannot send that game with the current bot",
	"GEO_POINT_INVALID":                   "Invalid geoposition provided",
	"GIF_CONTENT_TYPE_INVALID":            "GIF content-type invalid",
	"GIF_ID_INVALID":                      "The provided GIF ID is invalid",
	"GRAPH_INVALID_RELOAD":                "Invalid graph token provided, please reload the stats and provide the updated token",
	"GRAPH_OUTDATED_RELOAD":               "The graph is outdated, please get a new async token using stats.getBroadcastStats",
	"GROUPED_MEDIA_INVALID":               "Invalid grouped media",
	"HASH_INVALID":                        "The provided hash is invalid",
	"HISTORY_GET_FAILED":                  "Fetching of history failed",
	"IMAGE_PROCESS_FAILED":                "Failure while processing image",
	"INLINE_BOT_REQUIRED":                 "The action must be performed through an inline bot callback",
	"INLINE_RESULT_EXPIRED":               "The inline query expired",
	"INPUT_CONSTRUCTOR_INVALID":           "The provided constructor is invalid",
	"INPUT_FETCH_ERROR":                   "An error occurred while deserializing TL parameters",
	"INPUT_FETCH_FAIL":                    "Failed deserializing TL payload",
	"INPUT_LAYER_INVALID":                 "The provided layer is invalid",
	"INPUT_METHOD_INVALID":                "The invoked method does not exist anymore or has never existed",
	"INPUT_REQUEST_TOO_LONG":              "The request is too big",
	"INPUT_USER_DEACTIVATED":              "The specified user was deleted",
	"INVITE_HASH_EMPTY":                   "The invite hash is empty",
	"INVITE_HASH_EXPIRED":                 "The invite link has expired",
	"INVITE_HASH_INVALID":                 "The invite hash is invalid",
	"LANG_PACK_INVALID":                   "The provided language pack is invalid",
	"LASTNAME_INVALID":                    "The last name is invalid",
	"LIMIT_INVALID":                       "The provided limit is invalid",
	"LINK_NOT_MODIFIED":                   "Discussion link not modified",
	"LOCATION_INVALID":                    "The provided location is invalid",
	"MAX_ID_INVALID":                      "The provided max ID is invalid",
	"MAX_QTS_INVALID":                     "The provided QTS were invalid",
	"MD5_CHECKSUM_INVALID":                "The MD5 checksums do not match",
	"MEDIA_CAPTION_TOO_LONG":              "The caption is too long",
	"MEDIA_EMPTY":                         "The provided media object is invalid",
	"MEDIA_INVALID":                       "Media invalid",
	"MEDIA_NEW_INVALID":                   "The new media to edit the message with is invalid (such as stickers or voice notes)",
	"MEDIA_PREV_INVALID":                  "The old media cannot be edited with anything else (such as stickers or voice notes)",
	"MEGAGROUP_ID_INVALID":                "Invalid supergroup ID",
	"MEGAGROUP_PREHISTORY_HIDDEN":         "You can't set this discussion group because it's history is hidden",
	"MEGAGROUP_REQUIRED":                  "You can only use this method on a supergroup",
	"MEMBER_NO_LOCATION":                  "An internal failure occurred while fetching user info (couldn't find location)",
	"MEMBER_OCCUPY_PRIMARY_LOC_FAILED":    "Occupation of primary member location failed",
	"MESSAGE_AUTHOR_REQUIRED":             "Message author required",
	"MESSAGE_DELETE_FORBIDDEN":            "You can't delete one of the messages you tried to delete, most likely because it is a service message.",
	"MESSAGE_EDIT_TIME_EXPIRED":           "You can't edit this message anymore, too much time has passed since its creation.",
	"MESSAGE_EMPTY":                       "Empty or invalid UTF-8 message was sent",
	"MESSAGE_IDS_EMPTY":                   "No message ids were provided",
	"MESSAGE_ID_INVALID":                  "The provided message id is invalid",
	"MESSAGE_NOT_MODIFIED":                "Content of the message was not modified",
	"MESSAGE_POLL_CLOSED":                 "The poll was closed and can no longer be voted on",
	"MESSAGE_TOO_LONG":                    "Message was too long. Current maximum length is 4096 UTF-8 characters",
	"METHOD_INVALID":                      "The API method is invalid and cannot be used",
	"MSGID_DECREASE_RETRY":                "The request should be retried with a lower message ID",
	"MSG_ID_INVALID":                      "The message ID used in the peer was invalid",
	"MSG_WAIT_FAILED":                     "A waiting call returned an error",
	"MT_SEND_QUEUE_TOO_LONG":              "<DOESN'T HAVE ANY INFO ABOUT ERROR MT_SEND_QUEUE_TOO_LONG>",
	"MULTI_MEDIA_TOO_LONG":                "Too many media files for album",
	"NEED_CHAT_INVALID":                   "The provided chat is invalid",
	"NEED_MEMBER_INVALID":                 "The provided member is invalid or does not exist (for example a thumb size)",
	"NEW_SALT_INVALID":                    "The new salt is invalid",
	"NEW_SETTINGS_INVALID":                "The new settings are invalid",
	"OFFSET_INVALID":                      "The provided offset is invalid",
	"OFFSET_PEER_ID_INVALID":              "The provided offset peer is invalid",
	"OPTION_INVALID":                      "The option specified is invalid and does not exist in the target poll",
	"OPTIONS_TOO_MUCH":                    "You defined too many options for the poll",
	"PACK_SHORT_NAME_INVALID":             "Short pack name invalid",
	"PACK_SHORT_NAME_OCCUPIED":            "A stickerpack with this name already exists",
	"PACK_TITLE_INVALID":                  "The stickerpack title is invalid",
	"PARTICIPANT_CALL_FAILED":             "Failure while making call",
	"PARTICIPANT_VERSION_OUTDATED":        "The other participant does not use an up to date telegram client with support for calls",
	"PARTICIPANTS_TOO_FEW":                "Not enough participants",
	"PASSWORD_EMPTY":                      "The provided password is empty",
	"PASSWORD_HASH_INVALID":               "The password (and thus its hash value) you entered is invalid",
	"PASSWORD_MISSING":                    "You must enable 2FA in order to transfer ownership of a channel",
	"PASSWORD_REQUIRED":                   "The account must have 2-factor authentication enabled (a password) before this method can be used",
	"PAYMENT_PROVIDER_INVALID":            "The payment provider was not recognized or its token was invalid",
	"PEER_FLOOD":                          "Too many requests",
	"PEER_ID_INVALID":                     "The provided peer id is invalid",
	"PEER_ID_NOT_SUPPORTED":               "The provided peer ID is not supported",
	"PERSISTENT_TIMESTAMP_EMPTY":          "Persistent timestamp empty",
	"PERSISTENT_TIMESTAMP_INVALID":        "Persistent timestamp invalid",
	"PERSISTENT_TIMESTAMP_OUTDATED":       "Persistent timestamp outdated",
	"PHONE_CODE_EMPTY":                    "phone_code from a SMS is empty",
	"PHONE_CODE_EXPIRED":                  "The confirmation code has expired",
	"PHONE_CODE_HASH_EMPTY":               "The phone code hash is missing",
	"PHONE_CODE_INVALID":                  "Invalid SMS code was sent",
	"PHONE_NUMBER_APP_SIGNUP_FORBIDDEN":   "New accounts can be registrated only from official apps, this app doesn't allow it",
	"PHONE_NUMBER_BANNED":                 "The provided phone number is banned from telegram",
	"PHONE_NUMBER_FLOOD":                  "You asked for the code too many times.",
	"PHONE_NUMBER_INVALID":                "The phone number is invalid",
	"PHONE_NUMBER_OCCUPIED":               "The phone number is already in use",
	"PHONE_NUMBER_UNOCCUPIED":             "The code is valid but no user with the given number is registered",
	"PHONE_PASSWORD_FLOOD":                "You have tried logging in too many times",
	"PHONE_PASSWORD_PROTECTED":            "This phone is password protected",
	"PHOTO_CONTENT_TYPE_INVALID":          "Photo mime-type invalid",
	"PHOTO_CONTENT_URL_EMPTY":             "Photo URL invalid",
	"PHOTO_CROP_FILE_MISSING":             "Photo crop file missing",
	"PHOTO_CROP_SIZE_SMALL":               "Photo is too small",
	"PHOTO_EXT_INVALID":                   "The extension of the photo is invalid",
	"PHOTO_FILE_MISSING":                  "Profile photo file missing",
	"PHOTO_ID_INVALID":                    "Photo ID invalid",
	"PHOTO_INVALID":                       "Photo invalid",
	"PHOTO_INVALID_DIMENSIONS":            "The photo dimensions are invalid",
	"PHOTO_SAVE_FILE_INVALID":             "The photo you tried to send cannot be saved by Telegram. A reason may be that it exceeds 10MB. Try resizing it locally",
	"PHOTO_THUMB_URL_EMPTY":               "Photo thumbnail URL is empty",
	"PIN_RESTRICTED":                      "You can't pin messages in private chats with other people", // now you can, legacy error for lower api layers
	"PINNED_DIALOGS_TOO_MUCH":             "Too many pinned dialogs",
	"POLL_ANSWERS_INVALID":                "The poll did not have enough answers or had too many",
	"POLL_OPTION_DUPLICATE":               "A duplicate option was sent in the same poll",
	"POLL_OPTION_INVALID":                 "A poll option used invalid data (the data may be too long)",
	"POLL_QUESTION_INVALID":               "The poll question was either empty or too long",
	"POLL_UNSUPPORTED":                    "This layer does not support polls in the issued method",
	"POLL_VOTE_REQUIRED":                  "Cast a vote in the poll before calling this method",
	"PRIVACY_KEY_INVALID":                 "The privacy key is invalid",
	"PRIVACY_TOO_LONG":                    "Cannot add that many entities in a single request",
	"PRIVACY_VALUE_INVALID":               "The specified privacy rule combination is invalid",
	"PTS_CHANGE_EMPTY":                    "No PTS change",
	"QUERY_ID_EMPTY":                      "The query ID is empty",
	"QUERY_ID_INVALID":                    "The query ID is invalid",
	"QUERY_TOO_SHORT":                     "The query string is too short",
	"QUIZ_CORRECT_ANSWERS_EMPTY":          "A quiz must specify one correct answer",
	"QUIZ_CORRECT_ANSWERS_TOO_MUCH":       "There can only be one correct answer",
	"QUIZ_CORRECT_ANSWER_INVALID":         "An invalid value was provided to the correct_answers field",
	"QUIZ_MULTIPLE_INVALID":               "A poll cannot be both multiple choice and quiz",
	"RANDOM_ID_DUPLICATE":                 "You provided a random ID that was already used",
	"RANDOM_ID_EMPTY":                     "Random ID empty",
	"RANDOM_ID_INVALID":                   "A provided random ID is invalid",
	"RANDOM_LENGTH_INVALID":               "Random length invalid",
	"RANGES_INVALID":                      "Invalid range provided",
	"REACTION_EMPTY":                      "No reaction provided",
	"REACTION_INVALID":                    "Invalid reaction provided (only emoji are allowed)",
	"REG_ID_GENERATE_FAILED":              "Failure while generating registration ID",
	"REPLY_MARKUP_BUY_EMPTY":              "Reply markup for buy button empty",
	"REPLY_MARKUP_INVALID":                "The provided reply markup is invalid",
	"REPLY_MARKUP_TOO_LONG":               "The data embedded in the reply markup buttons was too much",
	"RESULT_ID_DUPLICATE":                 "Duplicated IDs on the sent results. Make sure to use unique IDs.",
	"RESULT_ID_EMPTY":                     "Result ID empty",
	"RESULT_TYPE_INVALID":                 "Result type invalid",
	"RESULTS_TOO_MUCH":                    "Too many results were provided",
	"REVOTE_NOT_ALLOWED":                  "You cannot change your vote",
	"RIGHT_FORBIDDEN":                     "Your admin rights do not allow you to do this",
	"RPC_CALL_FAIL":                       "Telegram is having internal issues, please try again later.",
	"RPC_MCGET_FAIL":                      "Telegram is having internal issues, please try again later.",
	"RSA_DECRYPT_FAILED":                  "Internal RSA decryption failed",
	"SCHEDULE_BOT_NOT_ALLOWED":            "Bots are not allowed to schedule messages",
	"SCHEDULE_DATE_INVALID":               "Invalid schedule date provided",
	"SCHEDULE_DATE_TOO_LATE":              "You can't schedule a message this far in the future",
	"SCHEDULE_STATUS_PRIVATE":             "You cannot schedule a message until the person comes online if their privacy does not show this information",
	"SCHEDULE_TOO_MUCH":                   "You cannot schedule more messages in this chat (last known limit of 100 per chat)",
	"SEARCH_QUERY_EMPTY":                  "The search query is empty",
	"SECONDS_INVALID":                     "Slow mode only supports certain values (0, 10s, 30s, 1m, 5m, 15m and 1h)",
	"SEND_MESSAGE_MEDIA_INVALID":          "The message media was invalid or not specified",
	"SEND_MESSAGE_TYPE_INVALID":           "The message type is invalid",
	"SESSION_EXPIRED":                     "The authorization has expired",
	"SESSION_PASSWORD_NEEDED":             "Two-steps verification is enabled and a password is required",
	"SESSION_REVOKED":                     "The authorization has been invalidated, because of the user terminating all sessions",
	"SETTINGS_INVALID":                    "Invalid settings were provided",
	"SHA256_HASH_INVALID":                 "The provided SHA256 hash is invalid",
	"SHORTNAME_OCCUPY_FAILED":             "An error occurred when trying to register the short-name used for the sticker pack. Try a different name",
	"SLOWMODE_MULTI_MSGS_DISABLED":        "Slowmode is enabled, you cannot forward multiple messages to this group.",
	"SMS_CODE_CREATE_FAILED":              "An error occurred while creating the SMS code",
	"SRP_ID_INVALID":                      "Invalid SRP ID provided",
	"SRP_PASSWORD_CHANGED":                "Password has changed",
	"START_PARAM_EMPTY":                   "The start parameter is empty",
	"START_PARAM_INVALID":                 "Start parameter invalid",
	"START_PARAM_TOO_LONG":                "Start parameter is too long",
	"STICKER_DOCUMENT_INVALID":            "The sticker file was invalid (this file has failed Telegram internal checks, make sure to use the correct format and comply with https://core.telegram.org/animated_stickers)",
	"STICKER_EMOJI_INVALID":               "Sticker emoji invalid",
	"STICKER_FILE_INVALID":                "Sticker file invalid",
	"STICKER_ID_INVALID":                  "The provided sticker ID is invalid",
	"STICKER_INVALID":                     "The provided sticker is invalid",
	"STICKER_PNG_DIMENSIONS":              "Sticker png dimensions invalid",
	"STICKER_PNG_NOPNG":                   "Stickers must be a png file but the used image was not a png",
	"STICKERS_EMPTY":                      "No sticker provided",
	"STICKERSET_INVALID":                  "The provided sticker set is invalid",
	"STORAGE_CHECK_FAILED":                "Server storage check failed",
	"STORE_INVALID_SCALAR_TYPE":           "<DOESN'T HAVE ANY INFO ABOUT ERROR STORE_INVALID_SCALAR_TYPE>",
	"TAKEOUT_INVALID":                     "The takeout session has been invalidated by another data export session",
	"TAKEOUT_REQUIRED":                    "A takeout session has to be initialized, first",
	"TEMP_AUTH_KEY_ALREADY_BOUND":         "The passed temporary key is already bound to another perm_auth_key_id",
	"TEMP_AUTH_KEY_EMPTY":                 "The request was not performed with a temporary authorization key",
	"THEME_FILE_INVALID":                  "Invalid theme file provided",
	"THEME_FORMAT_INVALID":                "Invalid theme format provided",
	"THEME_INVALID":                       "Invalid theme provided",
	"TMP_PASSWORD_DISABLED":               "The temporary password is disabled",
	"TOKEN_INVALID":                       "The provided token is invalid",
	"TTL_DAYS_INVALID":                    "The provided TTL is invalid",
	"TTL_MEDIA_INVALID":                   "Invalid media Time To Live was provided",
	"TYPE_CONSTRUCTOR_INVALID":            "The type constructor is invalid",
	"TYPES_EMPTY":                         "No top peer type was provided",
	"UNKNOWN_METHOD":                      "The method you tried to call cannot be called on non-CDN DCs",
	"UNTIL_DATE_INVALID":                  "That date cannot be specified in this request (try using nil)",
	"URL_INVALID":                         "The URL used was invalid (e.g. when answering a callback with an URL that's not t.me/yourbot or your game's URL)",
	"USERNAME_INVALID":                    `Nobody is using this username, or the username is unacceptable. If the latter, it must match ^[a-zA-Z][\w\d]{3,30}[a-zA-Z\d]&`,
	"USERNAME_NOT_MODIFIED":               "The username is not different from the current username",
	"USERNAME_NOT_OCCUPIED":               "The username is not in use by anyone else yet",
	"USERNAME_OCCUPIED":                   "The username is already taken",
	"USERS_TOO_FEW":                       "Not enough users (to create a chat, for example)",
	"USERS_TOO_MUCH":                      "The maximum number of users has been exceeded (to create a chat, for example)",
	"USER_ADMIN_INVALID":                  "Either you're not an admin or you tried to ban an admin that you didn't promote",
	"USER_ALREADY_PARTICIPANT":            "The authenticated user is already a participant of the chat",
	"USER_BANNED_IN_CHANNEL":              "You're banned from sending messages in supergroups/channels",
	"USER_BLOCKED":                        "User blocked",
	"USER_BOT":                            "Bots can only be admins in channels.",
	"USER_BOT_INVALID":                    "This method can only be called by a bot",
	"USER_BOT_REQUIRED":                   "This method can only be called by a bot",
	"USER_CHANNELS_TOO_MUCH":              "One of the users you tried to add is already in too many channels/supergroups",
	"USER_CREATOR":                        "You can't leave this channel, because you're its creator",
	"USER_DEACTIVATED_BAN":                "The user has been deactivated and banned",
	"USER_DEACTIVATED":                    "The user has been deleted/deactivated",
	"USER_ID_INVALID":                     "The provided user ID is invalid",
	"USER_INVALID":                        "The given user was invalid",
	"USER_IS_BLOCKED":                     "You were blocked by this user",
	"USER_IS_BOT":                         "Bots can't send messages to other bots",
	"USER_KICKED":                         "This user was kicked from this supergroup/channel",
	"USER_NOT_MUTUAL_CONTACT":             "The provided user is not a mutual contact",
	"USER_NOT_PARTICIPANT":                "The target user is not a member of the specified megagroup or channel",
	"USER_PRIVACY_RESTRICTED":             "The user's privacy settings do not allow you to do this",
	"USER_RESTRICTED":                     "You're spamreported, you can't create channels or chats.",
	"USERPIC_UPLOAD_REQUIRED":             "You must have a profile picture to publish your geolocation",
	"VIDEO_CONTENT_TYPE_INVALID":          "The video content type is not supported with the given parameters (i.e. supports_streaming)",
	"VIDEO_FILE_INVALID":                  "The specified video file is invalid",
	"WALLPAPER_FILE_INVALID":              "The given file cannot be used as a wallpaper",
	"WALLPAPER_INVALID":                   "The input wallpaper was not valid",
	"WC_CONVERT_URL_INVALID":              "WC convert URL invalid",
	"WEBDOCUMENT_INVALID":                 "Invalid webdocument URL provided",
	"WEBDOCUMENT_MIME_INVALID":            "Invalid webdocument mime type provided",
	"WEBDOCUMENT_SIZE_TOO_BIG":            "Webdocument is too big!",
	"WEBDOCUMENT_URL_INVALID":             "The given URL cannot be used",
	"WEBPAGE_CURL_FAILED":                 "Failure while fetching the webpage with cURL",
	"WEBPAGE_MEDIA_EMPTY":                 "Webpage media empty",
	"WORKER_BUSY_TOO_LONG_RETRY":          "Telegram workers are too busy to respond immediately",
	"YOU_BLOCKED_USER":                    "You blocked this user",

	// errors with additional data
	"2FA_CONFIRM_WAIT_X":        "You'll be able to reset your account in X seconds. If not, account will be deleted in 1 week for security reasons",
	"EMAIL_UNCONFIRMED_X":       "Email unconfirmed, the length of the code must be %v",
	"FILE_MIGRATE_X":            "The file to be accessed is currently stored in DC %v",
	"FILE_PART_X_MISSING":       "Part %v of the file is missing from storage",
	"FLOOD_TEST_PHONE_WAIT_X":   "A wait of %v seconds is required in the test servers",
	"FLOOD_WAIT_X":              "A wait of %v seconds is required",
	"INTERDC_X_CALL_ERROR":      "An error occurred while communicating with DC %v",
	"INTERDC_X_CALL_RICH_ERROR": "A rich error occurred while communicating with DC %v",
	"NETWORK_MIGRATE_X":         "The source IP address is associated with DC %v",
	"PASSWORD_TOO_FRESH_X":      "The password was added too recently and %v seconds must pass before using the method",
	"PHONE_MIGRATE_X":           "The phone number a user is trying to use for authorization is associated with DC %v",
	"SESSION_TOO_FRESH_X":       "The session logged in too recently and %v seconds must pass before calling the method",
	"SLOWMODE_WAIT_X":           "A wait of %v seconds is required before sending another message in this chat",
	"STATS_MIGRATE_X":           "The channel statistics must be fetched from DC %v",
	"TAKEOUT_INIT_DELAY_X":      "A wait of %v seconds is required before being able to initiate the takeout",
	"USER_MIGRATE_X":            "The user whose identity is being used to execute queries is associated with DC %v",

	// next errors was found, but they're too strange and looks like misspelling
	// "FILE_REFERENCE_*":                  "The file reference expired, it must be refreshed",
	//! pony floodwait https://core.telegram.org/method/messages.forwardMessages
	// "P0NY_FLOODWAIT":                    " ", //! No any description provided
	// "INPUT_METHOD_INVALID_1192227_X":    "Invalid method",
	// "INPUT_METHOD_INVALID_1400137063_X": "Invalid method",
	// "INPUT_METHOD_INVALID_1604042050_X": "Invalid method",
}

type BadMsgError struct {
	*objects.BadMsgNotification
	Description string
}

func BadMsgErrorFromNative(in *objects.BadMsgNotification) *BadMsgError {
	return &BadMsgError{
		BadMsgNotification: in,
		Description:        badMsgErrorCodes[uint8(in.Code)],
	}
}

func (e *BadMsgError) Error() string {
	return fmt.Sprintf("%v (code %v)", e.Description, e.Code)
}

// https://core.telegram.org/mtproto/service_messages_about_messages#notice-of-ignored-error-message
var badMsgErrorCodes = map[uint8]string{
	16: "msg_id too low (most likely, client time is wrong; it would be worthwhile to synchronize it using msg_id notifications and re-send the original message with the “correct” msg_id or wrap it in a container with a new msg_id if the original message had waited too long on the client to be transmitted)",
	17: "msg_id too high (similar to the previous case, the client time has to be synchronized, and the message re-sent with the correct msg_id",
	18: "incorrect two lower order msg_id bits (the server expects client message msg_id to be divisible by 4)",
	19: "container msg_id is the same as msg_id of a previously received message (this must never happen)",
	20: "message too old, and it cannot be verified whether the server has received a message with this msg_id or not",
	32: "msg_seqno too low (the server has already received a message with a lower msg_id but with either a higher or an equal and odd seqno)",
	33: "msg_seqno too high (similarly, there is a message with a higher msg_id but with either a lower or an equal and odd seqno)",
	34: "an even msg_seqno expected (irrelevant message), but odd received",
	35: "odd msg_seqno expected (relevant message), but even received",
	48: "incorrect server salt (in this case, the bad_server_salt response is received with the correct salt, and the message is to be re-sent with it)",
	64: "invalid container",
}

type BadSystemMessageCode int32

const (
	ErrBadMsgUnknown             BadSystemMessageCode = 0
	ErrBadMsgIdTooLow            BadSystemMessageCode = 16
	ErrBadMsgIdTooHigh           BadSystemMessageCode = 17
	ErrBadMsgIncorrectMsgIdBits  BadSystemMessageCode = 18
	ErrBadMsgWrongContainerMsgId BadSystemMessageCode = 19 // this must never happen
	ErrBadMsgMessageTooOld       BadSystemMessageCode = 20
	ErrBadMsgSeqNoTooLow         BadSystemMessageCode = 32
	ErrBadMsgSeqNoTooHigh        BadSystemMessageCode = 33
	ErrBadMsgSeqNoExpectedEven   BadSystemMessageCode = 34
	ErrBadMsgSeqNoExpectedOdd    BadSystemMessageCode = 35
	ErrBadMsgServerSaltIncorrect BadSystemMessageCode = 48
	ErrBadMsgInvalidContainer    BadSystemMessageCode = 64
)

// internal errors for internal purposes

type errorSessionConfigsChanged null

func (*errorSessionConfigsChanged) Error() string {
	return "session configuration was changed, need to repeat request"
}

func (*errorSessionConfigsChanged) CRC() uint32 {
	panic("makes no sense")
}
