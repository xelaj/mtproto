// uknown types, docs are not provided
// these items exist in tl schema, but in real life, they aren't using anywhere
// error#c4b9f9bb code:int text:string = Error;
// ipPort#d433ad73 ipv4:int port:int = IpPort;
// ipPortSecret#37982646 ipv4:int port:int secret:bytes = IpPort;
// accessPointRule#4679b65f phone_prefix_rules:string dc_id:int ips:Vector<IpPort> = AccessPointRule;
// help.configSimple#5a592a6c date:int expires:int rules:Vector<AccessPointRule> = help.ConfigSimple;

// @type Peer
// @constructor An empty constructor, no user or chat is defined.
inputPeerEmpty#7f3b18ea = InputPeer;
// @constructor Defines the current user.
inputPeerSelf#7da07ec9 = InputPeer;
// @constructor Defines a chat for further interaction.
// @param chat_id Chat idientifier
inputPeerChat#179be863 chat_id:int = InputPeer;
// @constructor Defines a user for further interaction.
// @param user_id User identifier
// @param access_hash **access_hash** value from the [user](https://core.telegram.org/constructor/user) constructor
inputPeerUser#7b8e7de6 user_id:int access_hash:long = InputPeer;
// @constructor Defines a channel for further interaction.
// @param channel_id Channel identifier
// @param access_hash **access_hash** value from the [channel](https://core.telegram.org/constructor/channel) constructor
inputPeerChannel#20adaef8 channel_id:int access_hash:long = InputPeer;
// @constructor Defines a [min](https://core.telegram.org/api/min) user that was seen in a certain message of a certain chat.
// @param peer The chat where the user was seen
// @param msg_id The message ID
// @param user_id The identifier of the user that was seen
inputPeerUserFromMessage#17bae2e6 peer:InputPeer msg_id:int user_id:int = InputPeer;
// @constructor Defines a [min](https://core.telegram.org/api/min) channel that was seen in a certain message of a certain chat.
// @param peer The chat where the channel's message was seen
// @param msg_id The message ID
// @param channel_id The identifier of the channel that was seen
inputPeerChannelFromMessage#9c95f7bb peer:InputPeer msg_id:int channel_id:int = InputPeer;

// @type Defines a user for subsequent interaction.
// @constructor Empty constructor, does not define a user.
inputUserEmpty#b98886cf = InputUser;
// @constructor Defines the current user.
inputUserSelf#f7c1b13f = InputUser;
// @constructor Defines a user for further interaction.
// @param access_hash access_hash value from the user constructor
// @param user_id     User identifier
inputUser#d8292816 user_id:int access_hash:long = InputUser;
// @constructor Defines a min user that was seen in a certain message of a certain chat.
// @param msg_id  The message ID
// @param peer    The chat where the user was seen
// @param user_id The identifier of the user that was seen
inputUserFromMessage#2d117597 peer:InputPeer msg_id:int user_id:int = InputUser;

// @type Object defines a contact from the user's phonebook.
// @constructor Phone contact. The client_id is just an arbitrary contact ID: it should be set, for example, to an incremental number when using contacts.importContacts, in order to retry importing only the contacts that weren't imported successfully.
// @param client_id  User identifier on the client
// @param first_name Contact's first name
// @param last_name  Contact's last name
// @param phone      Phone number
inputPhoneContact#f392b7f4 client_id:long phone:string first_name:string last_name:string = InputContact;

// @type Defines a file uploaded by the client.
// @constructor Defines a file saved in parts using the method upload.saveFilePart.
// @param id           Random file identifier created by the client
// @param md5_checksum In case the file's md5-hash was passed, contents of the file will be checked prior to use
// @param name         Full name of the file
// @param parts        Number of parts saved
inputFile#f52ff27f id:long parts:int name:string md5_checksum:string = InputFile;
// @constructor Assigns a big file (over 10Mb in size), saved in part using the method upload.saveBigFilePart.
// @param id    Random file id, created by the client
// @param name  Full file name
// @param parts Number of parts saved
inputFileBig#fa4f0bb5 id:long parts:int name:string = InputFile;

// @type Defines media content of a message.
// @constructor Empty media content of a message.
inputMediaEmpty#9664f57f = InputMedia;
// @constructor Photo
// @param file        The uploaded file
// @param stickers    Attached mask stickers
// @param ttl_seconds Time to live in seconds of self-destructing photo
inputMediaUploadedPhoto#1e287d04 flags:# file:InputFile stickers:flags.0?Vector<InputDocument> ttl_seconds:flags.1?int = InputMedia;
// @constructor
// @param id          Photo to be forwarded
// @param ttl_seconds Time to live in seconds of self-destructing photo
inputMediaPhoto#b3ba0635 flags:# id:InputPhoto ttl_seconds:flags.0?int = InputMedia;
// @constructor Map.
// @param geo_point GeoPoint
inputMediaGeoPoint#f9c44144 geo_point:InputGeoPoint = InputMedia;
// @constructor Phonebook contact
// @param first_name   Contact's first name
// @param last_name    Contact's last name
// @param phone_number Number of parts saved
// @param vcard        Contact vcard
inputMediaContact#f8ab7dfb phone_number:string first_name:string last_name:string vcard:string = InputMedia;
// @constructor New document
// @param nosound_video Whether the specified document is a video file with no audio tracks (a GIF animation (even as MPEG4), for example)
// @param force_file    Force the media file to be uploaded as document
// @param file          The uploaded file
// @param thumb         Thumbnail of the document, uploaded as for the file
// @param mime_type     MIME type of document
// @param attributes    Attributes that specify the type of the document (video, audio, voice, sticker, etc.)
// @param stickers      Attached stickers
// @param ttl_seconds   Time to live in seconds of self-destructing document
inputMediaUploadedDocument#5b38c6c1 flags:# nosound_video:flags.3?true force_file:flags.4?true file:InputFile thumb:flags.2?InputFile mime_type:string attributes:Vector<DocumentAttribute> stickers:flags.0?Vector<InputDocument> ttl_seconds:flags.1?int = InputMedia;
// @constructor Forwarded document
// @param id          The document to be forwarded.
// @param ttl_seconds Time to live of self-destructing document
inputMediaDocument#23ab23d2 flags:# id:InputDocument ttl_seconds:flags.0?int = InputMedia;
// @constructor Can be used to send a venue geolocation.
// @param geo_point  Geolocation
// @param title      Venue name
// @param address    Physical address of the venue
// @param provider   Venue provider: currently only "foursquare" needs to be supported
// @param venue_id   Venue ID in the provider's database
// @param venue_type Venue type in the provider's database
inputMediaVenue#c13d1c11 geo_point:InputGeoPoint title:string address:string provider:string venue_id:string venue_type:string = InputMedia;
// @constructor New photo that will be uploaded by the server using the specified URL
// @param url         URL of the photo"
// @param ttl_seconds Self-destruct time to live of photo"
inputMediaPhotoExternal#e5bbfe1a flags:# url:string ttl_seconds:flags.0?int = InputMedia;
// @constructor Document that will be downloaded by the telegram servers
// @param ttl_seconds Self-destruct time to live of document
// @param url         URL of the document
inputMediaDocumentExternal#fb52dc99 flags:# url:string ttl_seconds:flags.0?int = InputMedia;
// @constructor A game
// @param id The game to forward
inputMediaGame#d33f43f3 id:InputGame = InputMedia;
// @constructor Generated invoice of a bot payment
// @param title         roduct name, 1-32 characters
// @param description   Product description, 1-255 characters
// @param photo         URL of the product photo for the invoice. Can be a photo of the goods or a marketing image for a service. People like it better when they see what they are paying for.
// @param invoice       The actual invoice
// @param payload       Bot-defined invoice payload, 1-128 bytes. This will not be displayed to the user, use for your internal processes.
// @param provider      Payments provider token, obtained via Botfather
// @param provider_data JSON-encoded data about the invoice, which will be shared with the payment provider. A detailed description of required fields should be provided by the payment provider.
// @param start_param   Start parameter
inputMediaInvoice#f4e096c3 flags:# title:string description:string photo:flags.0?InputWebDocument invoice:Invoice payload:bytes provider:string provider_data:DataJSON start_param:string = InputMedia;
// @constructor Live geolocation
// @param stopped                       Whether sending of the geolocation was stopped
// @param geo_point                     Current geolocation
// @param heading                       For live locations, a direction in which the location moves, in degrees; 1-360.
// @param period                        Validity period of the current location
// @param proximity_notification_radius For live locations, a maximum distance to another chat member for proximity alerts, in meters (0-100000)
inputMediaGeoLive#971fa843 flags:# stopped:flags.0?true geo_point:InputGeoPoint heading:flags.2?int period:flags.1?int proximity_notification_radius:flags.3?int = InputMedia;
// @constructor A poll
// @param poll              The poll to send
// @param correct_answers   Correct answer IDs (for quiz polls)
// @param solution          Explanation of quiz solution
// @param solution_entities Message entities for styled text
inputMediaPoll#f94e5f1 flags:# poll:Poll correct_answers:flags.0?Vector<bytes> solution:flags.1?string solution_entities:flags.1?Vector<MessageEntity> = InputMedia;
// @constructor Assigns a big file (over 10Mb in size), saved in part using the method upload.saveBigFilePart.
// @param emoticon The emoji, for now 🏀, 🎲 and 🎯 are supported
inputMediaDice#e66fbf7b emoticon:string = InputMedia;

// @type Defines a new group profile photo.
// @constructor Empty constructor, remove group photo.
inputChatPhotoEmpty#1ca48f57 = InputChatPhoto;
// @constructor New photo to be set as group profile photo.
// @param file           File saved in parts using the method upload.saveFilePart
// @param video          Square video for animated profile picture"
// @param video_start_ts Timestamp that should be shown as static preview to the user (seconds)
inputChatUploadedPhoto#c642724e flags:# file:flags.0?InputFile video:flags.1?InputFile video_start_ts:flags.2?double = InputChatPhoto;
// @constructor Existing photo to be set as a chat profile photo.
// @param id Existing photo
inputChatPhoto#8953ad37 id:InputPhoto = InputChatPhoto;

// @type Defines a GeoPoint.
// @constructor Empty GeoPoint constructor.
inputGeoPointEmpty#e4c123d6 = InputGeoPoint;
// @constructor Defines a GeoPoint by its coordinates.
// @param accuracy_radius The estimated horizontal accuracy of the location, in meters; as defined by the sender.
// @param lat             Latitide
// @param long            Longtitude
inputGeoPoint#48222faf flags:# lat:double long:double accuracy_radius:flags.0?int = InputGeoPoint;

inputPhotoEmpty#1cd7bf0d = InputPhoto;
inputPhoto#3bb3b94a id:long access_hash:long file_reference:bytes = InputPhoto;

inputFileLocation#dfdaabe1 volume_id:long local_id:int secret:long file_reference:bytes = InputFileLocation;
inputEncryptedFileLocation#f5235d55 id:long access_hash:long = InputFileLocation;
inputDocumentFileLocation#bad07584 id:long access_hash:long file_reference:bytes thumb_size:string = InputFileLocation;
inputSecureFileLocation#cbc7ee28 id:long access_hash:long = InputFileLocation;
inputTakeoutFileLocation#29be5899 = InputFileLocation;
inputPhotoFileLocation#40181ffe id:long access_hash:long file_reference:bytes thumb_size:string = InputFileLocation;
inputPhotoLegacyFileLocation#d83466f3 id:long access_hash:long file_reference:bytes volume_id:long local_id:int secret:long = InputFileLocation;
inputPeerPhotoFileLocation#27d69997 flags:# big:flags.0?true peer:InputPeer volume_id:long local_id:int = InputFileLocation;
inputStickerSetThumb#dbaeae9 stickerset:InputStickerSet volume_id:long local_id:int = InputFileLocation;

// @type Chat partner or group.
// @constructor Chat partner
// @param user_id User identifier
peerUser#9db1bc6d user_id:int = Peer;
// @constructor Group.
// @param chat_id Group identifier
peerChat#bad0e5bb chat_id:int = Peer;
// @constructor Channel/supergroup
// @param channel_id Channel ID
peerChannel#bddde532 channel_id:int = Peer;

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

// @type Object defines a user.
// @constructor Empty constructor, non-existent user.
// @param id User identifier or `0`
userEmpty#200250ba id:int = User;
// @constructor Indicates info about a certain user
// @param self                   Whether this user indicates the currently logged in user
// @param contact                Whether this user is a contact
// @param mutual_contact         Whether this user is a mutual contact
// @param deleted                Whether the account of this user was deleted
// @param bot                    Is this user a bot?
// @param bot_chat_history       Can the bot see all messages in groups?
// @param bot_nochats            Can the bot be added to groups?
// @param verified               Whether this user is verified
// @param restricted             Access to this user must be restricted for the reason specified in restriction_reason
// @param min                    see [min](https://core.telegram.org/api/min)
// @param bot_inline_geo         Whether the bot can request our geolocation in inline mode
// @param support                Whether this is an official support user
// @param scam                   This may be a scam user
// @param apply_min_photo        If set, the profile picture for this user should be refetched
// @param id                     ID of the user
// @param access_hash            Access hash of the user
// @param first_name             First name
// @param last_name              Last name
// @param username               Username
// @param phone                  Phone number
// @param photo                  Profile picture of user
// @param status                 Online status of user
// @param bot_info_version       Version of the bot_info field in userFull, incremented every time it changes
// @param restriction_reason     Contains the reason why access to this user must be restricted.
// @param bot_inline_placeholder Inline placeholder for this inline bot
// @param lang_code              Language code of the user
user#938458c1 flags:# self:flags.10?true contact:flags.11?true mutual_contact:flags.12?true deleted:flags.13?true bot:flags.14?true bot_chat_history:flags.15?true bot_nochats:flags.16?true verified:flags.17?true restricted:flags.18?true min:flags.20?true bot_inline_geo:flags.21?true support:flags.23?true scam:flags.24?true apply_min_photo:flags.25?true id:int access_hash:flags.0?long first_name:flags.1?string last_name:flags.2?string username:flags.3?string phone:flags.4?string photo:flags.5?UserProfilePhoto status:flags.6?UserStatus bot_info_version:flags.14?int restriction_reason:flags.18?Vector<RestrictionReason> bot_inline_placeholder:flags.19?string lang_code:flags.22?string = User;

userProfilePhotoEmpty#4f11bae1 = UserProfilePhoto;
userProfilePhoto#69d3ab26 flags:# has_video:flags.0?true photo_id:long photo_small:FileLocation photo_big:FileLocation dc_id:int = UserProfilePhoto;

userStatusEmpty#9d05049 = UserStatus;
userStatusOnline#edb93949 expires:int = UserStatus;
userStatusOffline#8c703f was_online:int = UserStatus;
userStatusRecently#e26f42f1 = UserStatus;
userStatusLastWeek#7bf09fc = UserStatus;
userStatusLastMonth#77ebc742 = UserStatus;

// @type Object defines a group.
// @constructor
// @param id Group identifier
chatEmpty#9ba2d800 id:int = Chat;
// @constructor Info about a group
// @param creator               Whether the current user is the creator of the group
// @param kicked                Whether the current user was kicked from the group
// @param left                  Whether the current user has left the group
// @param deactivated           Whether the group was migrated
// @param id                    ID of the group
// @param title                 Title
// @param photo                 Chat photo
// @param participants_count    Participant count
// @param date                  Date of creation of the group
// @param version               Used in basic groups to reorder updates and make sure that all of them were received.
// @param migrated_to           Means this chat was upgraded to a supergroup
// @param admin_rights          Admin rights of the user in the group
// @param default_banned_rights Default banned rights of all users in the group
chat#3bda1bde flags:# creator:flags.0?true kicked:flags.1?true left:flags.2?true deactivated:flags.5?true id:int title:string photo:ChatPhoto participants_count:int date:int version:int migrated_to:flags.6?InputChannel admin_rights:flags.14?ChatAdminRights default_banned_rights:flags.18?ChatBannedRights = Chat;
// @constructor A group to which the user has no access. E.g., because the user was kicked from the group.
// @param id    User identifier
// @param title Group name
chatForbidden#7328bdb id:int title:string = Chat;
// @constructor Channel/supergroup info
// @param creator               Whether the current user is the creator of this channel
// @param left                  Whether the current user has left this channel
// @param broadcast             Is this a channel?
// @param verified              Is this channel verified by telegram?
// @param megagroup             Is this a supergroup?
// @param restricted            Whether viewing/writing in this channel for a reason (see restriction_reason
// @param signatures            Whether signatures are enabled (channels)
// @param min                   See min
// @param scam                  This channel/supergroup is probably a scam
// @param has_link              Whether this channel has a private join link
// @param has_geo               Whether this chanel has a geoposition
// @param slowmode_enabled      Whether slow mode is enabled for groups to prevent flood in chat
// @param id                    ID of the channel
// @param access_hash           Access hash
// @param title                 Title
// @param username              Username
// @param photo                 Profile photo
// @param date                  Date when the user joined the supergroup/channel, or if the user isn't a member, its creation date
// @param version               Version of the channel (always 0)
// @param restriction_reason    Contains the reason why access to this channel must be restricted.
// @param admin_rights          Admin rights of the user in this channel (see rights)
// @param banned_rights         Banned rights of the user in this channel (see rights)
// @param default_banned_rights Default chat rights (see rights)
// @param participants_count    Participant count
channel#d31a961e flags:# creator:flags.0?true left:flags.2?true broadcast:flags.5?true verified:flags.7?true megagroup:flags.8?true restricted:flags.9?true signatures:flags.11?true min:flags.12?true scam:flags.19?true has_link:flags.20?true has_geo:flags.21?true slowmode_enabled:flags.22?true call_active:flags.23?true id:int access_hash:flags.13?long title:string username:flags.6?string photo:ChatPhoto date:int version:int restriction_reason:flags.9?Vector<RestrictionReason> admin_rights:flags.14?ChatAdminRights banned_rights:flags.15?ChatBannedRights default_banned_rights:flags.18?ChatBannedRights participants_count:flags.17?int = Chat;
// @constructor Indicates a channel/supergroup we can't access because we were banned, or for some other reason.
// @param broadcast   Is this a channel
// @param megagroup   Is this a supergroup
// @param id          Channel ID
// @param access_hash Access hash
// @param title       Title
// @param until_date  The ban is valid until the specified date
channelForbidden#289da732 flags:# broadcast:flags.5?true megagroup:flags.8?true id:int access_hash:long title:string until_date:flags.16?int = Chat;

// @type Object containing detailed group info
// @constructor Detailed chat info
// @param about            About string for this chat
// @param bot_info         Info about bots that are in this chat
// @param can_set_username Can we change the username of this chat
// @param chat_photo       Chat photo
// @param exported_invite  Chat invite
// @param folder_id        Peer folder ID, for more info click here
// @param has_scheduled    Whether scheduled messages are available
// @param id               ID of the chat
// @param notify_settings  Notification settings
// @param participants     Participant list
// @param pinned_msg_id    Message ID of the last pinned message
chatFull#1b7c9db3 flags:# can_set_username:flags.7?true has_scheduled:flags.8?true id:int about:string participants:ChatParticipants chat_photo:flags.2?Photo notify_settings:PeerNotifySettings exported_invite:ExportedChatInvite bot_info:flags.3?Vector<BotInfo> pinned_msg_id:flags.6?int folder_id:flags.11?int = ChatFull;
// @constructor Full info about a channel/supergroup
// @param about                   Info about the channel
// @param admins_count            Number of channel admins
// @param available_min_id        Identifier of a maximum unavailable message in a channel due to hidden history.
// @param banned_count            Number of users banned from the channel
// @param blocked                 Whether any anonymous admin of this supergroup was blocked: if set, you won't receive messages from anonymous group admins in discussion replies via @replies
// @param bot_info                Info about bots in the channel/supergrup
// @param can_set_location        Can we set the geolocation of this group (for geogroups)
// @param can_set_stickers        Can we associate a stickerpack to the supergroup?
// @param can_set_username        Can we set the channel's username?
// @param can_view_participants   Can we vew the participant list?
// @param can_view_stats          Can the user view channel/supergroup statistics
// @param chat_photo              Channel picture
// @param exported_invite         Invite link
// @param folder_id               Peer folder ID, for more info click here
// @param has_scheduled           Whether scheduled messages are available
// @param hidden_prehistory       Is the history before we joined hidden to us?
// @param id                      ID of the channel
// @param kicked_count            Number of users kicked from the channel
// @param linked_chat_id          ID of the linked discussion chat for channels
// @param location                Location of the geogroup
// @param migrated_from_chat_id   The chat ID from which this group was migrated
// @param migrated_from_max_id    The message ID in the original chat at which this group was migrated
// @param notify_settings         Notification settings
// @param online_count            Number of users currently online
// @param participants_count      Number of participants of the channel
// @param pinned_msg_id           Message ID of the last pinned message
// @param pts                     Latest PTS for this channel
// @param read_inbox_max_id       Position up to which all incoming messages are read.
// @param read_outbox_max_id      Position up to which all outgoing messages are read.
// @param slowmode_next_send_date Indicates when the user will be allowed to send another message in the supergroup (unixdate)
// @param slowmode_seconds        If specified, users in supergroups will only be able to send one message every slowmode_seconds seconds
// @param stats_dc                If set, specifies the DC to use for fetching channel statistics
// @param stickerset              Associated stickerset
// @param unread_count            Count of unread messages
channelFull#f0e6672a flags:# can_view_participants:flags.3?true can_set_username:flags.6?true can_set_stickers:flags.7?true hidden_prehistory:flags.10?true can_set_location:flags.16?true has_scheduled:flags.19?true can_view_stats:flags.20?true blocked:flags.22?true id:int about:string participants_count:flags.0?int admins_count:flags.1?int kicked_count:flags.2?int banned_count:flags.2?int online_count:flags.13?int read_inbox_max_id:int read_outbox_max_id:int unread_count:int chat_photo:Photo notify_settings:PeerNotifySettings exported_invite:ExportedChatInvite bot_info:Vector<BotInfo> migrated_from_chat_id:flags.4?int migrated_from_max_id:flags.4?int pinned_msg_id:flags.5?int stickerset:flags.8?StickerSet available_min_id:flags.9?int folder_id:flags.11?int linked_chat_id:flags.14?int location:flags.15?ChannelLocation slowmode_seconds:flags.17?int slowmode_next_send_date:flags.18?int stats_dc:flags.12?int pts:int = ChatFull;

chatParticipant#c8d7493e user_id:int inviter_id:int date:int = ChatParticipant;
chatParticipantCreator#da13538a user_id:int = ChatParticipant;
chatParticipantAdmin#e2d6e436 user_id:int inviter_id:int date:int = ChatParticipant;

chatParticipantsForbidden#fc900c2b flags:# chat_id:int self_participant:flags.0?ChatParticipant = ChatParticipants;
chatParticipants#3f460fed chat_id:int participants:Vector<ChatParticipant> version:int = ChatParticipants;

chatPhotoEmpty#37c1011c = ChatPhoto;
chatPhoto#d20b9f3c flags:# has_video:flags.0?true photo_small:FileLocation photo_big:FileLocation dc_id:int = ChatPhoto;

// @type Object describing a message.
// @constructor Empty constructor, non-existent message.
// @param Message identifier
messageEmpty#83e5de54 id:int = Message;
// @constructor A message
// @param out                Is this an outgoing message
// @param mentioned          Whether we were mentioned in this message
// @param media_unread       Whether there are unread media attachments in this message
// @param silent             Whether this is a silent message (no notification triggered)
// @param post               Whether this is a channel post
// @param from_scheduled     Whether this is a scheduled message
// @param legacy             This is a legacy message: it has to be refetched with the new layer
// @param edit_hide          Whether the message should be shown as not modified to the user, even if an edit date is present
// @param pinned             Whether this message is pinned
// @param id                 ID of the message
// @param from_id            ID of the sender of the message
// @param peer_id            Peer ID, the chat where this message was sent
// @param fwd_from           Info about forwarded messages
// @param via_bot_id         ID of the inline bot that generated the message
// @param reply_to           Reply information
// @param date               Date of the message
// @param message            The message
// @param media              Media attachment
// @param reply_markup       Reply markup (bot/inline keyboards)
// @param entities           Message entities for styled text
// @param views              View count for channel posts
// @param forwards           Forward counter
// @param replies            Info about post comments (for channels) or message replies (for groups)
// @param edit_date          Last edit date of this message
// @param post_author        Name of the author of this message for channel posts (with signatures enabled)
// @param grouped_id         Multiple media messages sent using messages.sendMultiMedia with the same grouped ID indicate an album or media group
// @param restriction_reason Contains the reason why access to this message must be restricted.
message#58ae39c9 flags:# out:flags.1?true mentioned:flags.4?true media_unread:flags.5?true silent:flags.13?true post:flags.14?true from_scheduled:flags.18?true legacy:flags.19?true edit_hide:flags.21?true pinned:flags.24?true id:int from_id:flags.8?Peer peer_id:Peer fwd_from:flags.2?MessageFwdHeader via_bot_id:flags.11?int reply_to:flags.3?MessageReplyHeader date:int message:string media:flags.9?MessageMedia reply_markup:flags.6?ReplyMarkup entities:flags.7?Vector<MessageEntity> views:flags.10?int forwards:flags.10?int replies:flags.23?MessageReplies edit_date:flags.15?int post_author:flags.16?string grouped_id:flags.17?long restriction_reason:flags.22?Vector<RestrictionReason> = Message;
// @constructor Indicates a service message
// @param out          Whether the message is outgoing
// @param mentioned    Whether we were mentioned in the message
// @param media_unread Whether the message contains unread media
// @param silent       Whether the message is silent
// @param post         Whether it's a channel post
// @param legacy       This is a legacy message: it has to be refetched with the new layer
// @param id           Message ID
// @param from_id      ID of the sender of this message
// @param peer_id      Sender of service message
// @param reply_to     Reply (thread) information
// @param date         Message date
// @param action       Event connected with the service message
messageService#286fa604 flags:# out:flags.1?true mentioned:flags.4?true media_unread:flags.5?true silent:flags.13?true post:flags.14?true legacy:flags.19?true id:int from_id:flags.8?Peer peer_id:Peer reply_to:flags.3?MessageReplyHeader date:int action:MessageAction = Message;

messageMediaEmpty#3ded6320 = MessageMedia;
messageMediaPhoto#695150d7 flags:# photo:flags.0?Photo ttl_seconds:flags.2?int = MessageMedia;
messageMediaGeo#56e0d474 geo:GeoPoint = MessageMedia;
messageMediaContact#cbf24940 phone_number:string first_name:string last_name:string vcard:string user_id:int = MessageMedia;
messageMediaUnsupported#9f84f49e = MessageMedia;
messageMediaDocument#9cb070d7 flags:# document:flags.0?Document ttl_seconds:flags.2?int = MessageMedia;
messageMediaWebPage#a32dd600 webpage:WebPage = MessageMedia;
messageMediaVenue#2ec0533f geo:GeoPoint title:string address:string provider:string venue_id:string venue_type:string = MessageMedia;
messageMediaGame#fdb19008 game:Game = MessageMedia;
messageMediaInvoice#84551347 flags:# shipping_address_requested:flags.1?true test:flags.3?true title:string description:string photo:flags.0?WebDocument receipt_msg_id:flags.2?int currency:string total_amount:long start_param:string = MessageMedia;
messageMediaGeoLive#b940c666 flags:# geo:GeoPoint heading:flags.0?int period:int proximity_notification_radius:flags.1?int = MessageMedia;
messageMediaPoll#4bd6e798 poll:Poll results:PollResults = MessageMedia;
messageMediaDice#3f7ee58b value:int emoticon:string = MessageMedia;

messageActionEmpty#b6aef7b0 = MessageAction;
messageActionChatCreate#a6638b9a title:string users:Vector<int> = MessageAction;
messageActionChatEditTitle#b5a1ce5a title:string = MessageAction;
messageActionChatEditPhoto#7fcb13a8 photo:Photo = MessageAction;
messageActionChatDeletePhoto#95e3fbef = MessageAction;
messageActionChatAddUser#488a7337 users:Vector<int> = MessageAction;
messageActionChatDeleteUser#b2ae9b0c user_id:int = MessageAction;
messageActionChatJoinedByLink#f89cf5e8 inviter_id:int = MessageAction;
messageActionChannelCreate#95d2ac92 title:string = MessageAction;
messageActionChatMigrateTo#51bdb021 channel_id:int = MessageAction;
messageActionChannelMigrateFrom#b055eaee title:string chat_id:int = MessageAction;
messageActionPinMessage#94bd38ed = MessageAction;
messageActionHistoryClear#9fbab604 = MessageAction;
messageActionGameScore#92a72876 game_id:long score:int = MessageAction;
messageActionPaymentSentMe#8f31b327 flags:# currency:string total_amount:long payload:bytes info:flags.0?PaymentRequestedInfo shipping_option_id:flags.1?string charge:PaymentCharge = MessageAction;
messageActionPaymentSent#40699cd0 currency:string total_amount:long = MessageAction;
messageActionPhoneCall#80e11a7f flags:# video:flags.2?true call_id:long reason:flags.0?PhoneCallDiscardReason duration:flags.1?int = MessageAction;
messageActionScreenshotTaken#4792929b = MessageAction;
messageActionCustomAction#fae69f56 message:string = MessageAction;
messageActionBotAllowed#abe9affe domain:string = MessageAction;
messageActionSecureValuesSentMe#1b287353 values:Vector<SecureValue> credentials:SecureCredentialsEncrypted = MessageAction;
messageActionSecureValuesSent#d95c6154 types:Vector<SecureValueType> = MessageAction;
messageActionContactSignUp#f3f25f76 = MessageAction;
messageActionGeoProximityReached#98e0d697 from_id:Peer to_id:Peer distance:int = MessageAction;

dialog#2c171f72 flags:# pinned:flags.2?true unread_mark:flags.3?true peer:Peer top_message:int read_inbox_max_id:int read_outbox_max_id:int unread_count:int unread_mentions_count:int notify_settings:PeerNotifySettings pts:flags.0?int draft:flags.1?DraftMessage folder_id:flags.4?int = Dialog;
dialogFolder#71bd134c flags:# pinned:flags.2?true folder:Folder peer:Peer top_message:int unread_muted_peers_count:int unread_unmuted_peers_count:int unread_muted_messages_count:int unread_unmuted_messages_count:int = Dialog;

photoEmpty#2331b22d id:long = Photo;
photo#fb197a65 flags:# has_stickers:flags.0?true id:long access_hash:long file_reference:bytes date:int sizes:Vector<PhotoSize> video_sizes:flags.1?Vector<VideoSize> dc_id:int = Photo;

photoSizeEmpty#e17e23c type:string = PhotoSize;
photoSize#77bfb61b type:string location:FileLocation w:int h:int size:int = PhotoSize;
photoCachedSize#e9a734fa type:string location:FileLocation w:int h:int bytes:bytes = PhotoSize;
photoStrippedSize#e0b0bc2e type:string bytes:bytes = PhotoSize;
photoSizeProgressive#5aa86a51 type:string location:FileLocation w:int h:int sizes:Vector<int> = PhotoSize;
photoPathSize#d8214d41 type:string bytes:bytes = PhotoSize;

geoPointEmpty#1117dd5f = GeoPoint;
geoPoint#b2a2f663 flags:# long:double lat:double access_hash:long accuracy_radius:flags.0?int = GeoPoint;

auth.sentCode#5e002502 flags:# type:auth.SentCodeType phone_code_hash:string next_type:flags.1?auth.CodeType timeout:flags.2?int = auth.SentCode;

auth.authorization#cd050916 flags:# tmp_sessions:flags.0?int user:User = auth.Authorization;
auth.authorizationSignUpRequired#44747e9a flags:# terms_of_service:flags.0?help.TermsOfService = auth.Authorization;

auth.exportedAuthorization#df969c2d id:int bytes:bytes = auth.ExportedAuthorization;

inputNotifyPeer#b8bc5b0c peer:InputPeer = InputNotifyPeer;
inputNotifyUsers#193b4417 = InputNotifyPeer;
inputNotifyChats#4a95e84e = InputNotifyPeer;
inputNotifyBroadcasts#b1db7c7e = InputNotifyPeer;

inputPeerNotifySettings#9c3d198e flags:# show_previews:flags.0?Bool silent:flags.1?Bool mute_until:flags.2?int sound:flags.3?string = InputPeerNotifySettings;

peerNotifySettings#af509d20 flags:# show_previews:flags.0?Bool silent:flags.1?Bool mute_until:flags.2?int sound:flags.3?string = PeerNotifySettings;

peerSettings#733f2961 flags:# report_spam:flags.0?true add_contact:flags.1?true block_contact:flags.2?true share_contact:flags.3?true need_contacts_exception:flags.4?true report_geo:flags.5?true autoarchived:flags.7?true geo_distance:flags.6?int = PeerSettings;

wallPaper#a437c3ed id:long flags:# creator:flags.0?true default:flags.1?true pattern:flags.3?true dark:flags.4?true access_hash:long slug:string document:Document settings:flags.2?WallPaperSettings = WallPaper;
wallPaperNoFile#8af40b25 flags:# default:flags.1?true dark:flags.4?true settings:flags.2?WallPaperSettings = WallPaper;

inputReportReasonSpam#58dbcab8 = ReportReason;
inputReportReasonViolence#1e22c78d = ReportReason;
inputReportReasonPornography#2e59d922 = ReportReason;
inputReportReasonChildAbuse#adf44ee3 = ReportReason;
inputReportReasonOther#e1746d0a text:string = ReportReason;
inputReportReasonCopyright#9b89f93a = ReportReason;
inputReportReasonGeoIrrelevant#dbd4feed = ReportReason;

userFull#edf17c12 flags:# blocked:flags.0?true phone_calls_available:flags.4?true phone_calls_private:flags.5?true can_pin_message:flags.7?true has_scheduled:flags.12?true video_calls_available:flags.13?true user:User about:flags.1?string settings:PeerSettings profile_photo:flags.2?Photo notify_settings:PeerNotifySettings bot_info:flags.3?BotInfo pinned_msg_id:flags.6?int common_chats_count:int folder_id:flags.11?int = UserFull;

contact#f911c994 user_id:int mutual:Bool = Contact;

importedContact#d0028438 user_id:int client_id:long = ImportedContact;

contactStatus#d3680c61 user_id:int status:UserStatus = ContactStatus;

// @type Info on the current user's contact list.
// @constructor Contact list on the server is the same as the list on the client.
contacts.contactsNotModified#b74ba9d2 = contacts.Contacts;
// @constructor The current user's contact list and info on users.
// @param contacts    Contact list
// @param saved_count Number of contacts that were saved successfully
// @param users       User list
contacts.contacts#eae87e42 contacts:Vector<Contact> saved_count:int users:Vector<User> = contacts.Contacts;

contacts.importedContacts#77d01c3b imported:Vector<ImportedContact> popular_invites:Vector<PopularContact> retry_contacts:Vector<long> users:Vector<User> = contacts.ImportedContacts;

contacts.blocked#ade1591 blocked:Vector<PeerBlocked> chats:Vector<Chat> users:Vector<User> = contacts.Blocked;
contacts.blockedSlice#e1664194 count:int blocked:Vector<PeerBlocked> chats:Vector<Chat> users:Vector<User> = contacts.Blocked;

messages.dialogs#15ba6c40 dialogs:Vector<Dialog> messages:Vector<Message> chats:Vector<Chat> users:Vector<User> = messages.Dialogs;
messages.dialogsSlice#71e094f3 count:int dialogs:Vector<Dialog> messages:Vector<Message> chats:Vector<Chat> users:Vector<User> = messages.Dialogs;
messages.dialogsNotModified#f0e3e596 count:int = messages.Dialogs;

messages.messages#8c718e87 messages:Vector<Message> chats:Vector<Chat> users:Vector<User> = messages.Messages;
messages.messagesSlice#3a54685e flags:# inexact:flags.1?true count:int next_rate:flags.0?int offset_id_offset:flags.2?int messages:Vector<Message> chats:Vector<Chat> users:Vector<User> = messages.Messages;
messages.channelMessages#64479808 flags:# inexact:flags.1?true pts:int count:int offset_id_offset:flags.2?int messages:Vector<Message> chats:Vector<Chat> users:Vector<User> = messages.Messages;
messages.messagesNotModified#74535f21 count:int = messages.Messages;

messages.chats#64ff9fd5 chats:Vector<Chat> = messages.Chats;
messages.chatsSlice#9cd81144 count:int chats:Vector<Chat> = messages.Chats;

messages.chatFull#e5d7d19c full_chat:ChatFull chats:Vector<Chat> users:Vector<User> = messages.ChatFull;

messages.affectedHistory#b45c69d1 pts:int pts_count:int offset:int = messages.AffectedHistory;

// @type Object describes message filter.
// @constructor Filter is absent.
inputMessagesFilterEmpty#57e2f66c = MessagesFilter;
// @constructor Filter for messages containing photos.
inputMessagesFilterPhotos#9609a51c = MessagesFilter;
// @constructor Filter for messages containing videos.
inputMessagesFilterVideo#9fc00e65 = MessagesFilter;
// @constructor Filter for messages containing photos or videos.
inputMessagesFilterPhotoVideo#56e9f0e4 = MessagesFilter;
// @constructor Filter for messages containing documents.
inputMessagesFilterDocument#9eddf188 = MessagesFilter;
// @constructor Return only messages containing URLs
inputMessagesFilterUrl#7ef0dd87 = MessagesFilter;
// @constructor Return only messages containing gifs
inputMessagesFilterGif#ffc86587 = MessagesFilter;
// @constructor Return only messages containing voice notes
inputMessagesFilterVoice#50f5c392 = MessagesFilter;
// @constructor Return only messages containing audio files
inputMessagesFilterMusic#3751b49e = MessagesFilter;
// @constructor Return only chat photo changes
inputMessagesFilterChatPhotos#3a20ecb8 = MessagesFilter;
// @constructor Return only phone calls
// @param missed Return only missed phone calls
inputMessagesFilterPhoneCalls#80c99768 flags:# missed:flags.0?true = MessagesFilter;
// @constructor Return only round videos and voice notes
inputMessagesFilterRoundVoice#7a7c17a4 = MessagesFilter;
// @constructor Return only round videos
inputMessagesFilterRoundVideo#b549da53 = MessagesFilter;
// @constructor Return only messages where the current user was mentioned
inputMessagesFilterMyMentions#c1f8e69a = MessagesFilter;
// @constructor Return only messages containing geolocations
inputMessagesFilterGeo#e7026d0d = MessagesFilter;
// @constructor Return only messages containing contacts
inputMessagesFilterContacts#e062db83 = MessagesFilter;
// @constructor Fetch only pinned messages
inputMessagesFilterPinned#1bb00451 = MessagesFilter;

// @type Object contains info on events occured.
// @constructor New message.
updateNewMessage#1f2b0afd message:Message pts:int pts_count:int = Update;
// @constructor Sent message with random_id client identifier was assigned an identifier.
updateMessageID#4e90bfd6 id:int random_id:long = Update;
// @constructor Messages were deleted.
updateDeleteMessages#a20db0e5 messages:Vector<int> pts:int pts_count:int = Update;
// @constructor The user is preparing a message; typing, recording, uploading, etc. This update is valid for 6 seconds. If no repeated update received after 6 seconds, it should be considered that the user stopped doing whatever he's been doing.
updateUserTyping#5c486927 user_id:int action:SendMessageAction = Update;
// @constructor The user is preparing a message in a group; typing, recording, uploading, etc. This update is valid for 6 seconds. If no repeated update received after 6 seconds, it should be considered that the user stopped doing whatever he's been doing.
updateChatUserTyping#9a65ea1f chat_id:int user_id:int action:SendMessageAction = Update;
// @constructor Composition of chat participants changed.
updateChatParticipants#7761198 participants:ChatParticipants = Update;
// @constructor Contact status update.
updateUserStatus#1bfbd823 user_id:int status:UserStatus = Update;
// @constructor Changes the user's first name, last name and username.
updateUserName#a7332b73 user_id:int first_name:string last_name:string username:string = Update;
// @constructor Change of contact's profile photo.
updateUserPhoto#95313b0c user_id:int date:int photo:UserProfilePhoto previous:Bool = Update;
// @constructor New encrypted message.
updateNewEncryptedMessage#12bcbd9a message:EncryptedMessage qts:int = Update;
// @constructor Interlocutor is typing a message in an encrypted chat. Update period is 6 second. If upon this time there is no repeated update, it shall be considered that the interlocutor stopped typing.
updateEncryptedChatTyping#1710f156 chat_id:int = Update;
// @constructor Change of state in an encrypted chat.
updateEncryption#b4a2e88d chat:EncryptedChat date:int = Update;
// @constructor Communication history in an encrypted chat was marked as read.
updateEncryptedMessagesRead#38fe25b7 chat_id:int max_date:int date:int = Update;
// @constructor New group member.
updateChatParticipantAdd#ea4b0e5c chat_id:int user_id:int inviter_id:int date:int version:int = Update;
// @constructor A member has left the group.
updateChatParticipantDelete#6e5f8c22 chat_id:int user_id:int version:int = Update;
// @constructor Changes in the data center configuration options.
updateDcOptions#8e5e9873 dc_options:Vector<DcOption> = Update;
// @constructor Changes in notification settings.
updateNotifySettings#bec268ef peer:NotifyPeer notify_settings:PeerNotifySettings = Update;
// @constructor A service message for the user. The app must show the message to the user upon receiving this update. In case the popup parameter was passed, the text message must be displayed in a popup alert immediately upon receipt. It is recommended to handle the text as you would an ordinary message in terms of highlighting links, etc. The message must also be stored locally as part of the message history with the user id 777000 (Telegram Notifications).
updateServiceNotification#ebe46819 flags:# popup:flags.0?true inbox_date:flags.1?int type:string message:string media:MessageMedia entities:Vector<MessageEntity> = Update;
// @constructor Privacy rules were changed
updatePrivacy#ee3b272a key:PrivacyKey rules:Vector<PrivacyRule> = Update;
// @constructor A user's phone number was changed
updateUserPhone#12b9417b user_id:int phone:string = Update;
// @constructor Incoming messages were read
updateReadHistoryInbox#9c974fdf flags:# folder_id:flags.0?int peer:Peer max_id:int still_unread_count:int pts:int pts_count:int = Update;
// @constructor Outgoing messages were read
updateReadHistoryOutbox#2f2f21bf peer:Peer max_id:int pts:int pts_count:int = Update;
// @constructor An (IV) webpage preview was generated
updateWebPage#7f891213 webpage:WebPage pts:int pts_count:int = Update;
// @constructor Contents of messages in the common message box were read
updateReadMessagesContents#68c13933 messages:Vector<int> pts:int pts_count:int = Update;
// @constructor There are new updates in the specified channel, the client must fetch them, eventually starting the specified pts if the difference is too long or if the channel isn't currently in the states.
updateChannelTooLong#eb0467fb flags:# channel_id:int pts:flags.0?int = Update;
// @constructor A new channel is available
updateChannel#b6d45656 channel_id:int = Update;
// @constructor A new message was sent in a channel/supergroup
updateNewChannelMessage#62ba04d9 message:Message pts:int pts_count:int = Update;
// @constructor Incoming messages in a channel/supergroup were read
updateReadChannelInbox#330b5424 flags:# folder_id:flags.0?int channel_id:int max_id:int still_unread_count:int pts:int = Update;
// @constructor Some messages in a supergroup/channel were deleted
updateDeleteChannelMessages#c37521c9 channel_id:int messages:Vector<int> pts:int pts_count:int = Update;
// @constructor The view counter of a message in a channel has changed
updateChannelMessageViews#98a12b4b channel_id:int id:int views:int = Update;
// @constructor Admin permissions of a user in a legacy group were changed
updateChatParticipantAdmin#b6901959 chat_id:int user_id:int is_admin:Bool version:int = Update;
// @constructor A new stickerset was installed
updateNewStickerSet#688a30aa stickerset:messages.StickerSet = Update;
// @constructor The order of stickersets was changed
updateStickerSetsOrder#bb2d201 flags:# masks:flags.0?true order:Vector<long> = Update;
// @constructor Installed stickersets have changed, the client should refetch them using messages.getAllStickers
updateStickerSets#43ae3dec = Update;
// @constructor The saved gif list has changed, the client should refetch it using messages.getSavedGifs
updateSavedGifs#9375341e = Update;
// @constructor An incoming inline query
updateBotInlineQuery#54826690 flags:# query_id:long user_id:int query:string geo:flags.0?GeoPoint offset:string = Update;
// @constructor The result of an inline query that was chosen by a user and sent to their chat partner. Please see our documentation on the feedback collecting for details on how to enable these updates for your bot.
updateBotInlineSend#e48f964 flags:# user_id:int query:string geo:flags.0?GeoPoint id:string msg_id:flags.1?InputBotInlineMessageID = Update;
// @constructor A message was edited in a channel/supergroup
updateEditChannelMessage#1b3f4df7 message:Message pts:int pts_count:int = Update;
// @constructor A callback button was pressed, and the button data was sent to the bot that created the button
updateBotCallbackQuery#e73547e1 flags:# query_id:long user_id:int peer:Peer msg_id:int chat_instance:long data:flags.0?bytes game_short_name:flags.1?string = Update;
// @constructor A message was edited
updateEditMessage#e40370a3 message:Message pts:int pts_count:int = Update;
// @constructor This notification is received by bots when a button is pressed
updateInlineBotCallbackQuery#f9d27a5a flags:# query_id:long user_id:int msg_id:InputBotInlineMessageID chat_instance:long data:flags.0?bytes game_short_name:flags.1?string = Update;
// @constructor Outgoing messages in a channel/supergroup were read
updateReadChannelOutbox#25d6c9c7 channel_id:int max_id:int = Update;
// @constructor Notifies a change of a message draft.
updateDraftMessage#ee2bb969 peer:Peer draft:DraftMessage = Update;
// @constructor Some featured stickers were marked as read
updateReadFeaturedStickers#571d2742 = Update;
// @constructor The recent sticker list was updated
updateRecentStickers#9a422c20 = Update;
// @constructor The server-side configuration has changed; the client should re-fetch the config using help.getConfig
updateConfig#a229dd06 = Update;
// @constructor Common message box sequence PTS has changed, state has to be refetched using updates.getState
updatePtsChanged#3354678f = Update;
// @constructor A webpage preview of a link in a channel/supergroup message was generated
updateChannelWebPage#40771900 channel_id:int webpage:WebPage pts:int pts_count:int = Update;
// @constructor A dialog was pinned/unpinned
updateDialogPinned#6e6fe51c flags:# pinned:flags.0?true folder_id:flags.1?int peer:DialogPeer = Update;
// @constructor Pinned dialogs were updated
updatePinnedDialogs#fa0f3ca2 flags:# folder_id:flags.1?int order:flags.0?Vector<DialogPeer> = Update;
// @constructor A new incoming event; for bots only
updateBotWebhookJSON#8317c0c3 data:DataJSON = Update;
// @constructor A new incoming query; for bots only
updateBotWebhookJSONQuery#9b9240a6 query_id:long data:DataJSON timeout:int = Update;
// @constructor This object contains information about an incoming shipping query.
updateBotShippingQuery#e0cdc940 query_id:long user_id:int payload:bytes shipping_address:PostAddress = Update;
// @constructor This object contains information about an incoming pre-checkout query.
updateBotPrecheckoutQuery#5d2f3aa9 flags:# query_id:long user_id:int payload:bytes info:flags.0?PaymentRequestedInfo shipping_option_id:flags.1?string currency:string total_amount:long = Update;
// @constructor An incoming phone call
updatePhoneCall#ab0f6b1e phone_call:PhoneCall = Update;
// @constructor A language pack has changed, the client should manually fetch the changed strings using langpack.getDifference
updateLangPackTooLong#46560264 lang_code:string = Update;
// @constructor Language pack updated
updateLangPack#56022f4d difference:LangPackDifference = Update;
// @constructor The list of favorited stickers was changed, the client should call messages.getFavedStickers to refetch the new list
updateFavedStickers#e511996d = Update;
// @constructor The specified channel/supergroup messages were read
updateChannelReadMessagesContents#89893b45 channel_id:int messages:Vector<int> = Update;
// @constructor All contacts were deleted
updateContactsReset#7084a7be = Update;
// @constructor The history of a channel/supergroup was hidden.
updateChannelAvailableMessages#70db6837 channel_id:int available_min_id:int = Update;
// @constructor The manual unread mark of a chat was changed
updateDialogUnreadMark#e16459c3 flags:# unread:flags.0?true peer:DialogPeer = Update;
// @constructor The results of a poll have changed
updateMessagePoll#aca1657b flags:# poll_id:long poll:flags.0?Poll results:PollResults = Update;
// @constructor Default banned rights in a normal chat were updated
updateChatDefaultBannedRights#54c01850 peer:Peer default_banned_rights:ChatBannedRights version:int = Update;
// @constructor The dialog list of a folder was changed
updateFolderPeers#19360dc0 folder_peers:Vector<FolderPeer> pts:int pts_count:int = Update;
// @constructor Settings of a certain peer have changed
updatePeerSettings#6a7e7366 peer:Peer settings:PeerSettings = Update;
// @constructor List of peers near you was updated
updatePeerLocated#b4afcfb0 peers:Vector<PeerLocated> = Update;
// @constructor New incoming scheduled message
updateNewScheduledMessage#39a51dfb message:Message = Update;
// @constructor Some scheduled messages were deleted
updateDeleteScheduledMessages#90866cee peer:Peer messages:Vector<int> = Update;
// @constructor A new cloud theme was installed
updateTheme#8216fba3 theme:Theme = Update;
// next are already described
// @constructor Live geoposition message was viewed
// @param peer   The user that viewed the live geoposition
// @param msg_id Message ID of geoposition message
updateGeoLiveViewed#871fb939 peer:Peer msg_id:int = Update;
// @constructor A login token (for login via QR code) was generated
updateLoginToken#564fe691 = Update;
// @constructor A specific user has voted in a poll
// @param poll_id Poll ID
// @param user_id User ID
// @param options Chosen option(s)
updateMessagePollVote#42f88f2c poll_id:long user_id:int options:Vector<bytes> = Update;
// @constructor A new folder was added
// @param id     Folder ID
// @param filter Folder info
updateDialogFilter#26ffde7d flags:# id:int filter:flags.0?DialogFilter = Update;
// @constructor New chat folders order
// @param order Ordered folder IDs
updateDialogFilterOrder#a5d72105 order:Vector<int> = Update;
// @constructor Update folder list
updateDialogFilters#3504914f = Update;
// @constructor Incoming phone call signaling payload
// @param phone_call_id Phone call ID
// @param data          Signaling payload
updatePhoneCallSignalingData#2661bf09 phone_call_id:long data:bytes = Update;
// @constructor A participant has left/joined a channel or supergroup.
// @param channel_id       Channel ID
// @param date             Date of the event
// @param user_id          User in question
// @param prev_participant Previous participant status
// @param new_participant  New participant status
// @param qts              [PTS](https://core.telegram.org/api/updates)
updateChannelParticipant#65d2b464 flags:# channel_id:int date:int user_id:int prev_participant:flags.0?ChannelParticipant new_participant:flags.1?ChannelParticipant qts:int = Update;
// @constructor The forward counter of a message in a channel has changed
// @param id       ID of the message
// @param forwards New forward counter
updateChannelMessageForwards#6e8a84df channel_id:int id:int forwards:int = Update;
// @constructor Incoming messages were marked as read
// @param channel_id     Discussion group ID
// @param top_msg_id     ID of the group message that started the thread (message in linked discussion group)
// @param read_max_id    Message ID of latest read incoming message for this thread
// @param broadcast_id   If set, contains the ID of the channel that contains the post that started the comment thread in the discussion group (channel_id)
// @param broadcast_post If set, contains the ID of the channel post that started the the comment thread
updateReadChannelDiscussionInbox#1cc7de54 flags:# channel_id:int top_msg_id:int read_max_id:int broadcast_id:flags.0?int broadcast_post:flags.0?int = Update;
// @constructor Outgoing comments in a discussion thread were marked as read
// @param top_msg_id  ID of the group message that started the thread
// @param read_max_id Message ID of latest read outgoing message for this thread
updateReadChannelDiscussionOutbox#4638a26c channel_id:int top_msg_id:int read_max_id:int = Update;
// @constructor A peer was blocked
// @param blocked Whether the peer was blocked or unblocked
updatePeerBlocked#246a4b22 peer_id:Peer blocked:Bool = Update;
// @constructor user is typing in a supergroup, channel or message thread
// @param channel_id Channel ID
// @param top_msg_id Thread ID
// @param user_id    User ID
// @param action     Whether the user is typing, sending a media or doing something else
updateChannelUserTyping#ff2abe9f flags:# channel_id:int top_msg_id:flags.0?int user_id:int action:SendMessageAction = Update;
// @constructor Some messages were pinned in a chat
// @param pinned    Whether the messages were pinned or unpinned
// @param peer      Peer
// @param messages  Message IDs
// @param pts       Event count after generation
// @param pts_count Number of events that were generated
updatePinnedMessages#ed85eab5 flags:# pinned:flags.0?true peer:Peer messages:Vector<int> pts:int pts_count:int = Update;
// @constructor Messages were pinned/unpinned in a channel/supergroup
// @param pinned     Whether the messages were pinned or unpinned
// @param channel_id Channel ID
// @param messages   Messages
// @param pts        Event count after generation
// @param pts_count  Number of events that were generated
updatePinnedChannelMessages#8588878b flags:# pinned:flags.0?true channel_id:int messages:Vector<int> pts:int pts_count:int = Update;

updates.state#a56c2a3e pts:int qts:int date:int seq:int unread_count:int = updates.State;

updates.differenceEmpty#5d75a138 date:int seq:int = updates.Difference;
updates.difference#f49ca0 new_messages:Vector<Message> new_encrypted_messages:Vector<EncryptedMessage> other_updates:Vector<Update> chats:Vector<Chat> users:Vector<User> state:updates.State = updates.Difference;
updates.differenceSlice#a8fb1981 new_messages:Vector<Message> new_encrypted_messages:Vector<EncryptedMessage> other_updates:Vector<Update> chats:Vector<Chat> users:Vector<User> intermediate_state:updates.State = updates.Difference;
updates.differenceTooLong#4afe8f6d pts:int = updates.Difference;

updatesTooLong#e317af7e = Updates;
updateShortMessage#2296d2c8 flags:# out:flags.1?true mentioned:flags.4?true media_unread:flags.5?true silent:flags.13?true id:int user_id:int message:string pts:int pts_count:int date:int fwd_from:flags.2?MessageFwdHeader via_bot_id:flags.11?int reply_to:flags.3?MessageReplyHeader entities:flags.7?Vector<MessageEntity> = Updates;
updateShortChatMessage#402d5dbb flags:# out:flags.1?true mentioned:flags.4?true media_unread:flags.5?true silent:flags.13?true id:int from_id:int chat_id:int message:string pts:int pts_count:int date:int fwd_from:flags.2?MessageFwdHeader via_bot_id:flags.11?int reply_to:flags.3?MessageReplyHeader entities:flags.7?Vector<MessageEntity> = Updates;
updateShort#78d4dec1 update:Update date:int = Updates;
updatesCombined#725b04c3 updates:Vector<Update> users:Vector<User> chats:Vector<Chat> date:int seq_start:int seq:int = Updates;
updates#74ae4240 updates:Vector<Update> users:Vector<User> chats:Vector<Chat> date:int seq:int = Updates;
updateShortSentMessage#11f1331c flags:# out:flags.1?true id:int pts:int pts_count:int date:int media:flags.9?MessageMedia entities:flags.7?Vector<MessageEntity> = Updates;

photos.photos#8dca6aa5 photos:Vector<Photo> users:Vector<User> = photos.Photos;
photos.photosSlice#15051f54 count:int photos:Vector<Photo> users:Vector<User> = photos.Photos;

photos.photo#20212ca8 photo:Photo users:Vector<User> = photos.Photo;

upload.file#96a18d5 type:storage.FileType mtime:int bytes:bytes = upload.File;
upload.fileCdnRedirect#f18cda44 dc_id:int file_token:bytes encryption_key:bytes encryption_iv:bytes file_hashes:Vector<FileHash> = upload.File;

dcOption#18b7a10d flags:# ipv6:flags.0?true media_only:flags.1?true tcpo_only:flags.2?true cdn:flags.3?true static:flags.4?true id:int ip_address:string port:int secret:flags.10?bytes = DcOption;

// @type Object contains info on API configuring parameters.
// @constructor Current configuration
// @param phonecalls_enabled         Whether phone calls can be used
// @param default_p2p_contacts       Whether the client should use P2P by default for phone calls with contacts
// @param preload_featured_stickers  Whether the client should preload featured stickers
// @param ignore_phone_entities      Whether the client should ignore phone entities
// @param revoke_pm_inbox            Whether incoming private messages can be deleted for both participants
// @param blocked_mode               Indicates that telegram is probably censored by governments/ISPs in the current region
// @param pfs_enabled                Whether pfs was used
// @param date                       Current date at the server
// @param expires                    Expiration date of this config: when it expires it'll have to be refetched using HelpGetConfig
// @param test_mode                  Whether we're connected to the test DCs
// @param this_dc                    ID of the DC that returned the reply
// @param dc_options                 DC IP list
// @param dc_txt_domain_name         Domain name for fetching encrypted DC list from DNS TXT record
// @param chat_size_max              Maximum member count for normal groups
// @param megagroup_size_max         Maximum member count for supergroups
// @param forwarded_count_max        Maximum number of messages that can be forwarded at once using MessagesForwardMessages
// @param online_update_period_ms    The client should update its online status every N milliseconds
// @param offline_blur_timeout_ms    Delay before offline status needs to be sent to the server
// @param offline_idle_timeout_ms    Time without any user activity after which it should be treated offline
// @param online_cloud_timeout_ms    If we are offline, but were online from some other client in last online_cloud_timeout_ms milliseconds after we had gone offline, then delay offline notification for notify_cloud_delay_ms milliseconds.
// @param notify_cloud_delay_ms      If we are offline, but online from some other client then delay sending the offline notification for notify_cloud_delay_ms milliseconds.
// @param notify_default_delay_ms    If some other client is online, then delay notification for notification_default_delay_ms milliseconds
// @param push_chat_period_ms        Not for client use
// @param push_chat_limit            Not for client use
// @param saved_gifs_limit           Maximum count of saved gifs
// @param edit_time_limit            Only messages with age smaller than the one specified can be edited
// @param revoke_time_limit          Only channel/supergroup messages with age smaller than the specified can be deleted
// @param revoke_pm_time_limit       Only private messages with age smaller than the specified can be deleted
// @param rating_e_decay             Exponential decay rate for computing top peer rating
// @param stickers_recent_limit      Maximum number of recent stickers
// @param stickers_faved_limit       Maximum number of faved stickers
// @param channels_read_media_period Indicates that round videos (video notes) and voice messages sent in channels and older than the specified period must be marked as read
// @param tmp_sessions               Temporary passport sessions
// @param pinned_dialogs_count_max   Maximum count of pinned dialogs
// @param pinned_infolder_count_max  Maximum count of dialogs per folder
// @param call_receive_timeout_ms    Maximum allowed outgoing ring time in VoIP calls: if the user we're calling doesn't reply within the specified time (in milliseconds), we should hang up the call
// @param call_ring_timeout_ms       Maximum allowed incoming ring time in VoIP calls: if the current user doesn't reply within the specified time (in milliseconds), the call will be automatically refused
// @param call_connect_timeout_ms    VoIP connection timeout: if the instance of libtgvoip on the other side of the call doesn't connect to our instance of libtgvoip within the specified time (in milliseconds), the call must be aborted
// @param call_packet_timeout_ms     If during a VoIP call a packet isn't received for the specified period of time, the call must be aborted
// @param me_url_prefix              The domain to use to parse in-app links.For example t.me indicates that t.me/username links should parsed to @username, t.me/addsticker/name should be parsed to the appropriate stickerset and so on...
// @param autoupdate_url_prefix      URL to use to auto-update the current app
// @param gif_search_username        Username of the bot to use to search for GIFs
// @param venue_search_username      Username of the bot to use to search for venues
// @param img_search_username        Username of the bot to use for image search
// @param static_maps_provider       ID of the map provider to use for venues
// @param caption_length_max         Maximum length of caption (length in utf8 codepoints)
// @param message_length_max         Maximum length of messages (length in utf8 codepoints)
// @param webfile_dc_id              DC ID to use to download webfiles
// @param suggested_lang_code        Suggested language code
// @param lang_pack_version          Language pack version
// @param base_lang_pack_version     Basic language pack version
config#330b4067 flags:# phonecalls_enabled:flags.1?true default_p2p_contacts:flags.3?true preload_featured_stickers:flags.4?true ignore_phone_entities:flags.5?true revoke_pm_inbox:flags.6?true blocked_mode:flags.8?true pfs_enabled:flags.13?true date:int expires:int test_mode:Bool this_dc:int dc_options:Vector<DcOption> dc_txt_domain_name:string chat_size_max:int megagroup_size_max:int forwarded_count_max:int online_update_period_ms:int offline_blur_timeout_ms:int offline_idle_timeout_ms:int online_cloud_timeout_ms:int notify_cloud_delay_ms:int notify_default_delay_ms:int push_chat_period_ms:int push_chat_limit:int saved_gifs_limit:int edit_time_limit:int revoke_time_limit:int revoke_pm_time_limit:int rating_e_decay:int stickers_recent_limit:int stickers_faved_limit:int channels_read_media_period:int tmp_sessions:flags.0?int pinned_dialogs_count_max:int pinned_infolder_count_max:int call_receive_timeout_ms:int call_ring_timeout_ms:int call_connect_timeout_ms:int call_packet_timeout_ms:int me_url_prefix:string autoupdate_url_prefix:flags.7?string gif_search_username:flags.9?string venue_search_username:flags.10?string img_search_username:flags.11?string static_maps_provider:flags.12?string caption_length_max:int message_length_max:int webfile_dc_id:int suggested_lang_code:flags.2?string lang_pack_version:flags.2?int base_lang_pack_version:flags.2?int = Config;

nearestDc#8e1a1775 country:string this_dc:int nearest_dc:int = NearestDc;

help.appUpdate#1da7158f flags:# can_not_skip:flags.0?true id:int version:string text:string entities:Vector<MessageEntity> document:flags.1?Document url:flags.2?string = help.AppUpdate;
help.noAppUpdate#c45a6536 = help.AppUpdate;

help.inviteText#18cb9f78 message:string = help.InviteText;

encryptedChatEmpty#ab7ec0a0 id:int = EncryptedChat;
encryptedChatWaiting#3bf703dc id:int access_hash:long date:int admin_id:int participant_id:int = EncryptedChat;
encryptedChatRequested#62718a82 flags:# folder_id:flags.0?int id:int access_hash:long date:int admin_id:int participant_id:int g_a:bytes = EncryptedChat;
encryptedChat#fa56ce36 id:int access_hash:long date:int admin_id:int participant_id:int g_a_or_b:bytes key_fingerprint:long = EncryptedChat;
encryptedChatDiscarded#13d6dd27 id:int = EncryptedChat;

inputEncryptedChat#f141b5e1 chat_id:int access_hash:long = InputEncryptedChat;

encryptedFileEmpty#c21f497e = EncryptedFile;
encryptedFile#4a70994c id:long access_hash:long size:int dc_id:int key_fingerprint:int = EncryptedFile;

inputEncryptedFileEmpty#1837c364 = InputEncryptedFile;
inputEncryptedFileUploaded#64bd0306 id:long parts:int md5_checksum:string key_fingerprint:int = InputEncryptedFile;
inputEncryptedFile#5a17b5e5 id:long access_hash:long = InputEncryptedFile;
inputEncryptedFileBigUploaded#2dc173c8 id:long parts:int key_fingerprint:int = InputEncryptedFile;

encryptedMessage#ed18c118 random_id:long chat_id:int date:int bytes:bytes file:EncryptedFile = EncryptedMessage;
encryptedMessageService#23734b06 random_id:long chat_id:int date:int bytes:bytes = EncryptedMessage;

messages.dhConfigNotModified#c0e24635 random:bytes = messages.DhConfig;
messages.dhConfig#2c221edd g:int p:bytes version:int random:bytes = messages.DhConfig;

messages.sentEncryptedMessage#560f8935 date:int = messages.SentEncryptedMessage;
messages.sentEncryptedFile#9493ff32 date:int file:EncryptedFile = messages.SentEncryptedMessage;

inputDocumentEmpty#72f0eaae = InputDocument;
inputDocument#1abfb575 id:long access_hash:long file_reference:bytes = InputDocument;

documentEmpty#36f8c871 id:long = Document;
document#1e87342b flags:# id:long access_hash:long file_reference:bytes date:int mime_type:string size:int thumbs:flags.0?Vector<PhotoSize> video_thumbs:flags.1?Vector<VideoSize> dc_id:int attributes:Vector<DocumentAttribute> = Document;

help.support#17c6b5f6 phone_number:string user:User = help.Support;

notifyPeer#9fd40bd8 peer:Peer = NotifyPeer;
notifyUsers#b4c83b4c = NotifyPeer;
notifyChats#c007cec3 = NotifyPeer;
notifyBroadcasts#d612e8ef = NotifyPeer;

sendMessageTypingAction#16bf744e = SendMessageAction;
sendMessageCancelAction#fd5ec8f5 = SendMessageAction;
sendMessageRecordVideoAction#a187d66f = SendMessageAction;
sendMessageUploadVideoAction#e9763aec progress:int = SendMessageAction;
sendMessageRecordAudioAction#d52f73f7 = SendMessageAction;
sendMessageUploadAudioAction#f351d7ab progress:int = SendMessageAction;
sendMessageUploadPhotoAction#d1d34a26 progress:int = SendMessageAction;
sendMessageUploadDocumentAction#aa0cd9e4 progress:int = SendMessageAction;
sendMessageGeoLocationAction#176f8ba1 = SendMessageAction;
sendMessageChooseContactAction#628cbc6f = SendMessageAction;
sendMessageGamePlayAction#dd6a8f48 = SendMessageAction;
sendMessageRecordRoundAction#88f27fbc = SendMessageAction;
sendMessageUploadRoundAction#243e1c66 progress:int = SendMessageAction;

contacts.found#b3134d9d my_results:Vector<Peer> results:Vector<Peer> chats:Vector<Chat> users:Vector<User> = contacts.Found;

// @type Privacy key
// @enum Whether we can see the exact last online timestamp of the user
inputPrivacyKeyStatusTimestamp#4f96cb18 = InputPrivacyKey;
// @enum Whether the user can be invited to chats
inputPrivacyKeyChatInvite#bdfb0426 = InputPrivacyKey;
// @enum Whether the user will accept phone calls
inputPrivacyKeyPhoneCall#fabadc5f = InputPrivacyKey;
// @enum Whether the user allows P2P communication during VoIP calls
inputPrivacyKeyPhoneP2P#db9e70d2 = InputPrivacyKey;
// @enum Whether messages forwarded from this user will be [anonymous](https://telegram.org/blog/unsend-privacy-emoji#anonymous-forwarding)
inputPrivacyKeyForwards#a4dd4c08 = InputPrivacyKey;
// @enum Whether people will be able to see the user's profile picture
inputPrivacyKeyProfilePhoto#5719bacc = InputPrivacyKey;
// @enum Whether people will be able to see the user's phone number
inputPrivacyKeyPhoneNumber#352dafa = InputPrivacyKey;
// @enum Whether people can add you to their contact list by your phone number
inputPrivacyKeyAddedByPhone#d1219bdd = InputPrivacyKey;

// @type Privacy key
// @enum Whether we can see the last online timestamp
privacyKeyStatusTimestamp#bc2eab30 = PrivacyKey;
// @enum Whether the user can be invited to chats
privacyKeyChatInvite#500e6dfa = PrivacyKey;
// @enum Whether the user accepts phone calls
privacyKeyPhoneCall#3d662b7b = PrivacyKey;
// @enum Whether P2P connections in phone calls are allowed
privacyKeyPhoneP2P#39491cc8 = PrivacyKey;
// @enum Whether messages forwarded from the user will be anonymously forwarded
privacyKeyForwards#69ec56a3 = PrivacyKey;
// @enum Whether the profile picture of the user is visible
privacyKeyProfilePhoto#96151fed = PrivacyKey;
// @enum Whether the user allows us to see his phone number
privacyKeyPhoneNumber#d19ae46d = PrivacyKey;
// @enum Whether people can add you to their contact list by your phone number
privacyKeyAddedByPhone#42ffd42b = PrivacyKey;

inputPrivacyValueAllowContacts#d09e07b = InputPrivacyRule;
inputPrivacyValueAllowAll#184b35ce = InputPrivacyRule;
inputPrivacyValueAllowUsers#131cc67f users:Vector<InputUser> = InputPrivacyRule;
inputPrivacyValueDisallowContacts#ba52007 = InputPrivacyRule;
inputPrivacyValueDisallowAll#d66b66c9 = InputPrivacyRule;
inputPrivacyValueDisallowUsers#90110467 users:Vector<InputUser> = InputPrivacyRule;
inputPrivacyValueAllowChatParticipants#4c81c1ba chats:Vector<int> = InputPrivacyRule;
inputPrivacyValueDisallowChatParticipants#d82363af chats:Vector<int> = InputPrivacyRule;

privacyValueAllowContacts#fffe1bac = PrivacyRule;
privacyValueAllowAll#65427b82 = PrivacyRule;
privacyValueAllowUsers#4d5bbe0c users:Vector<int> = PrivacyRule;
privacyValueDisallowContacts#f888fa1a = PrivacyRule;
privacyValueDisallowAll#8b73e763 = PrivacyRule;
privacyValueDisallowUsers#c7f49b7 users:Vector<int> = PrivacyRule;
privacyValueAllowChatParticipants#18be796b chats:Vector<int> = PrivacyRule;
privacyValueDisallowChatParticipants#acae0690 chats:Vector<int> = PrivacyRule;

account.privacyRules#50a04e45 rules:Vector<PrivacyRule> chats:Vector<Chat> users:Vector<User> = account.PrivacyRules;

accountDaysTTL#b8d0afdf days:int = AccountDaysTTL;

documentAttributeImageSize#6c37c15c w:int h:int = DocumentAttribute;
documentAttributeAnimated#11b58939 = DocumentAttribute;
documentAttributeSticker#6319d612 flags:# mask:flags.1?true alt:string stickerset:InputStickerSet mask_coords:flags.0?MaskCoords = DocumentAttribute;
documentAttributeVideo#ef02ce6 flags:# round_message:flags.0?true supports_streaming:flags.1?true duration:int w:int h:int = DocumentAttribute;
documentAttributeAudio#9852f9c6 flags:# voice:flags.10?true duration:int title:flags.0?string performer:flags.1?string waveform:flags.2?bytes = DocumentAttribute;
documentAttributeFilename#15590068 file_name:string = DocumentAttribute;
documentAttributeHasStickers#9801d2f7 = DocumentAttribute;

messages.stickersNotModified#f1749a22 = messages.Stickers;
messages.stickers#e4599bbd hash:int stickers:Vector<Document> = messages.Stickers;

stickerPack#12b299d4 emoticon:string documents:Vector<long> = StickerPack;

messages.allStickersNotModified#e86602c3 = messages.AllStickers;
messages.allStickers#edfd405f hash:int sets:Vector<StickerSet> = messages.AllStickers;

messages.affectedMessages#84d19185 pts:int pts_count:int = messages.AffectedMessages;

// @type Instant View webpage preview
// @constructor No preview is available for the webpage
// @param Preview ID
webPageEmpty#eb1477e8 id:long = WebPage;
// @constructor A preview of the webpage is currently being generated
// @param date When was the processing started
// @param id   ID of preview
webPagePending#c586da1c id:long date:int = WebPage;
// @constructor Webpage preview
// @param attributes   Webpage attributes
// @param author       Author of the content
// @param cached_page  Page contents in instant view format
// @param description  Content description
// @param display_url  Webpage URL to be displayed to the user
// @param document     Preview of the content as a media file
// @param duration     Duration of the content, in seconds
// @param embed_height Height of the embedded preview
// @param embed_type   MIME type of the embedded preview, (e.g., text/html or video/mp4)
// @param embed_url    URL to show in the embedded preview
// @param embed_width  Width of the embedded preview
// @param flags        Flags, see TL conditional fields
// @param hash         Hash for pagination, for more info click here
// @param id           Preview ID
// @param photo        Image representing the content
// @param site_name    Short name of the site (e.g., Google Docs, App Store)
// @param title        Title of the content
// @param type         Type of the web page. Can be: article, photo, audio, video, document, profile, app, or something else
// @param url          URL of previewed webpage
webPage#e89c45b2 flags:# id:long url:string display_url:string hash:int type:flags.0?string site_name:flags.1?string title:flags.2?string description:flags.3?string photo:flags.4?Photo embed_url:flags.5?string embed_type:flags.5?string embed_width:flags.6?int embed_height:flags.6?int duration:flags.7?int author:flags.8?string document:flags.9?Document cached_page:flags.10?Page attributes:flags.12?Vector<WebPageAttribute> = WebPage;
// @constructor The preview of the webpage hasn't changed
// @param cached_page_views Page view count
// @param flags             Flags, see TL conditional fields
webPageNotModified#7311ca11 flags:# cached_page_views:flags.0?int = WebPage;

authorization#ad01d61d flags:# current:flags.0?true official_app:flags.1?true password_pending:flags.2?true hash:long device_model:string platform:string system_version:string api_id:int app_name:string app_version:string date_created:int date_active:int ip:string country:string region:string = Authorization;

account.authorizations#1250abde authorizations:Vector<Authorization> = account.Authorizations;

account.password#ad2641f8 flags:# has_recovery:flags.0?true has_secure_values:flags.1?true has_password:flags.2?true current_algo:flags.2?PasswordKdfAlgo srp_B:flags.2?bytes srp_id:flags.2?long hint:flags.3?string email_unconfirmed_pattern:flags.4?string new_algo:PasswordKdfAlgo new_secure_algo:SecurePasswordKdfAlgo secure_random:bytes = account.Password;

account.passwordSettings#9a5c33e5 flags:# email:flags.0?string secure_settings:flags.1?SecureSecretSettings = account.PasswordSettings;

account.passwordInputSettings#c23727c9 flags:# new_algo:flags.0?PasswordKdfAlgo new_password_hash:flags.0?bytes hint:flags.0?string email:flags.1?string new_secure_settings:flags.2?SecureSecretSettings = account.PasswordInputSettings;

auth.passwordRecovery#137948a5 email_pattern:string = auth.PasswordRecovery;

receivedNotifyMessage#a384b779 id:int flags:int = ReceivedNotifyMessage;

chatInviteEmpty#69df3769 = ExportedChatInvite;
chatInviteExported#fc2e05bc link:string = ExportedChatInvite;

chatInviteAlready#5a686d7c chat:Chat = ChatInvite;
chatInvite#dfc2f58e flags:# channel:flags.0?true broadcast:flags.1?true public:flags.2?true megagroup:flags.3?true title:string photo:Photo participants_count:int participants:flags.4?Vector<User> = ChatInvite;
chatInvitePeek#61695cb0 chat:Chat expires:int = ChatInvite;

inputStickerSetEmpty#ffb62b95 = InputStickerSet;
inputStickerSetID#9de7a269 id:long access_hash:long = InputStickerSet;
inputStickerSetShortName#861cc8a0 short_name:string = InputStickerSet;
inputStickerSetAnimatedEmoji#28703c8 = InputStickerSet;
inputStickerSetDice#e67f520e emoticon:string = InputStickerSet;

stickerSet#eeb46f27 flags:# archived:flags.1?true official:flags.2?true masks:flags.3?true animated:flags.5?true installed_date:flags.0?int id:long access_hash:long title:string short_name:string thumb:flags.4?PhotoSize thumb_dc_id:flags.4?int count:int hash:int = StickerSet;

messages.stickerSet#b60a24a6 set:StickerSet packs:Vector<StickerPack> documents:Vector<Document> = messages.StickerSet;

botCommand#c27ac8c7 command:string description:string = BotCommand;

botInfo#98e81d3a user_id:int description:string commands:Vector<BotCommand> = BotInfo;

keyboardButton#a2fa4880 text:string = KeyboardButton;
keyboardButtonUrl#258aff05 text:string url:string = KeyboardButton;
keyboardButtonCallback#35bbdb6b flags:# requires_password:flags.0?true text:string data:bytes = KeyboardButton;
keyboardButtonRequestPhone#b16a6c29 text:string = KeyboardButton;
keyboardButtonRequestGeoLocation#fc796b3f text:string = KeyboardButton;
keyboardButtonSwitchInline#568a748 flags:# same_peer:flags.0?true text:string query:string = KeyboardButton;
keyboardButtonGame#50f41ccf text:string = KeyboardButton;
keyboardButtonBuy#afd93fbb text:string = KeyboardButton;
keyboardButtonUrlAuth#10b78d29 flags:# text:string fwd_text:flags.0?string url:string button_id:int = KeyboardButton;
inputKeyboardButtonUrlAuth#d02e7fd4 flags:# request_write_access:flags.0?true text:string fwd_text:flags.1?string url:string bot:InputUser = KeyboardButton;
keyboardButtonRequestPoll#bbc7515d flags:# quiz:flags.0?Bool text:string = KeyboardButton;

keyboardButtonRow#77608b83 buttons:Vector<KeyboardButton> = KeyboardButtonRow;

replyKeyboardHide#a03e5b85 flags:# selective:flags.2?true = ReplyMarkup;
replyKeyboardForceReply#f4108aa0 flags:# single_use:flags.1?true selective:flags.2?true = ReplyMarkup;
replyKeyboardMarkup#3502758c flags:# resize:flags.0?true single_use:flags.1?true selective:flags.2?true rows:Vector<KeyboardButtonRow> = ReplyMarkup;
replyInlineMarkup#48a30254 rows:Vector<KeyboardButtonRow> = ReplyMarkup;

messageEntityUnknown#bb92ba95 offset:int length:int = MessageEntity;
messageEntityMention#fa04579d offset:int length:int = MessageEntity;
messageEntityHashtag#6f635b0d offset:int length:int = MessageEntity;
messageEntityBotCommand#6cef8ac7 offset:int length:int = MessageEntity;
messageEntityUrl#6ed02538 offset:int length:int = MessageEntity;
messageEntityEmail#64e475c2 offset:int length:int = MessageEntity;
messageEntityBold#bd610bc9 offset:int length:int = MessageEntity;
messageEntityItalic#826f8b60 offset:int length:int = MessageEntity;
messageEntityCode#28a20571 offset:int length:int = MessageEntity;
messageEntityPre#73924be0 offset:int length:int language:string = MessageEntity;
messageEntityTextUrl#76a6d327 offset:int length:int url:string = MessageEntity;
messageEntityMentionName#352dca58 offset:int length:int user_id:int = MessageEntity;
inputMessageEntityMentionName#208e68c9 offset:int length:int user_id:InputUser = MessageEntity;
messageEntityPhone#9b69e34b offset:int length:int = MessageEntity;
messageEntityCashtag#4c4e743f offset:int length:int = MessageEntity;
messageEntityUnderline#9c4e7e8b offset:int length:int = MessageEntity;
messageEntityStrike#bf0693d4 offset:int length:int = MessageEntity;
messageEntityBlockquote#20df5d0 offset:int length:int = MessageEntity;
messageEntityBankCard#761e6af4 offset:int length:int = MessageEntity;

inputChannelEmpty#ee8c1e86 = InputChannel;
inputChannel#afeb712e channel_id:int access_hash:long = InputChannel;
inputChannelFromMessage#2a286531 peer:InputPeer msg_id:int channel_id:int = InputChannel;

contacts.resolvedPeer#7f077ad9 peer:Peer chats:Vector<Chat> users:Vector<User> = contacts.ResolvedPeer;

messageRange#ae30253 min_id:int max_id:int = MessageRange;

updates.channelDifferenceEmpty#3e11affb flags:# final:flags.0?true pts:int timeout:flags.1?int = updates.ChannelDifference;
updates.channelDifferenceTooLong#a4bcc6fe flags:# final:flags.0?true timeout:flags.1?int dialog:Dialog messages:Vector<Message> chats:Vector<Chat> users:Vector<User> = updates.ChannelDifference;
updates.channelDifference#2064674e flags:# final:flags.0?true pts:int timeout:flags.1?int new_messages:Vector<Message> other_updates:Vector<Update> chats:Vector<Chat> users:Vector<User> = updates.ChannelDifference;

channelMessagesFilterEmpty#94d42ee7 = ChannelMessagesFilter;
channelMessagesFilter#cd77d957 flags:# exclude_new_messages:flags.1?true ranges:Vector<MessageRange> = ChannelMessagesFilter;

channelParticipant#15ebac1d user_id:int date:int = ChannelParticipant;
channelParticipantSelf#a3289a6d user_id:int inviter_id:int date:int = ChannelParticipant;
channelParticipantCreator#447dca4b flags:# user_id:int admin_rights:ChatAdminRights rank:flags.0?string = ChannelParticipant;
channelParticipantAdmin#ccbebbaf flags:# can_edit:flags.0?true self:flags.1?true user_id:int inviter_id:flags.1?int promoted_by:int date:int admin_rights:ChatAdminRights rank:flags.2?string = ChannelParticipant;
channelParticipantBanned#1c0facaf flags:# left:flags.0?true user_id:int kicked_by:int date:int banned_rights:ChatBannedRights = ChannelParticipant;
channelParticipantLeft#c3c6796b user_id:int = ChannelParticipant;

channelParticipantsRecent#de3f3c79 = ChannelParticipantsFilter;
channelParticipantsAdmins#b4608969 = ChannelParticipantsFilter;
channelParticipantsKicked#a3b54985 q:string = ChannelParticipantsFilter;
channelParticipantsBots#b0d1865b = ChannelParticipantsFilter;
channelParticipantsBanned#1427a5e1 q:string = ChannelParticipantsFilter;
channelParticipantsSearch#656ac4b q:string = ChannelParticipantsFilter;
channelParticipantsContacts#bb6ae88d q:string = ChannelParticipantsFilter;
channelParticipantsMentions#e04b5ceb flags:# q:flags.0?string top_msg_id:flags.1?int = ChannelParticipantsFilter;

channels.channelParticipants#f56ee2a8 count:int participants:Vector<ChannelParticipant> users:Vector<User> = channels.ChannelParticipants;
channels.channelParticipantsNotModified#f0173fe9 = channels.ChannelParticipants;

channels.channelParticipant#d0d9b163 participant:ChannelParticipant users:Vector<User> = channels.ChannelParticipant;

help.termsOfService#780a0310 flags:# popup:flags.0?true id:DataJSON text:string entities:Vector<MessageEntity> min_age_confirm:flags.1?int = help.TermsOfService;

messages.savedGifsNotModified#e8025ca2 = messages.SavedGifs;
messages.savedGifs#2e0709a5 hash:int gifs:Vector<Document> = messages.SavedGifs;

inputBotInlineMessageMediaAuto#3380c786 flags:# message:string entities:flags.1?Vector<MessageEntity> reply_markup:flags.2?ReplyMarkup = InputBotInlineMessage;
inputBotInlineMessageText#3dcd7a87 flags:# no_webpage:flags.0?true message:string entities:flags.1?Vector<MessageEntity> reply_markup:flags.2?ReplyMarkup = InputBotInlineMessage;
inputBotInlineMessageMediaGeo#96929a85 flags:# geo_point:InputGeoPoint heading:flags.0?int period:flags.1?int proximity_notification_radius:flags.3?int reply_markup:flags.2?ReplyMarkup = InputBotInlineMessage;
inputBotInlineMessageMediaVenue#417bbf11 flags:# geo_point:InputGeoPoint title:string address:string provider:string venue_id:string venue_type:string reply_markup:flags.2?ReplyMarkup = InputBotInlineMessage;
inputBotInlineMessageMediaContact#a6edbffd flags:# phone_number:string first_name:string last_name:string vcard:string reply_markup:flags.2?ReplyMarkup = InputBotInlineMessage;
inputBotInlineMessageGame#4b425864 flags:# reply_markup:flags.2?ReplyMarkup = InputBotInlineMessage;

inputBotInlineResult#88bf9319 flags:# id:string type:string title:flags.1?string description:flags.2?string url:flags.3?string thumb:flags.4?InputWebDocument content:flags.5?InputWebDocument send_message:InputBotInlineMessage = InputBotInlineResult;
inputBotInlineResultPhoto#a8d864a7 id:string type:string photo:InputPhoto send_message:InputBotInlineMessage = InputBotInlineResult;
inputBotInlineResultDocument#fff8fdc4 flags:# id:string type:string title:flags.1?string description:flags.2?string document:InputDocument send_message:InputBotInlineMessage = InputBotInlineResult;
inputBotInlineResultGame#4fa417f2 id:string short_name:string send_message:InputBotInlineMessage = InputBotInlineResult;

botInlineMessageMediaAuto#764cf810 flags:# message:string entities:flags.1?Vector<MessageEntity> reply_markup:flags.2?ReplyMarkup = BotInlineMessage;
botInlineMessageText#8c7f65e2 flags:# no_webpage:flags.0?true message:string entities:flags.1?Vector<MessageEntity> reply_markup:flags.2?ReplyMarkup = BotInlineMessage;
botInlineMessageMediaGeo#51846fd flags:# geo:GeoPoint heading:flags.0?int period:flags.1?int proximity_notification_radius:flags.3?int reply_markup:flags.2?ReplyMarkup = BotInlineMessage;
botInlineMessageMediaVenue#8a86659c flags:# geo:GeoPoint title:string address:string provider:string venue_id:string venue_type:string reply_markup:flags.2?ReplyMarkup = BotInlineMessage;
botInlineMessageMediaContact#18d1cdc2 flags:# phone_number:string first_name:string last_name:string vcard:string reply_markup:flags.2?ReplyMarkup = BotInlineMessage;

botInlineResult#11965f3a flags:# id:string type:string title:flags.1?string description:flags.2?string url:flags.3?string thumb:flags.4?WebDocument content:flags.5?WebDocument send_message:BotInlineMessage = BotInlineResult;
botInlineMediaResult#17db940b flags:# id:string type:string photo:flags.0?Photo document:flags.1?Document title:flags.2?string description:flags.3?string send_message:BotInlineMessage = BotInlineResult;

messages.botResults#947ca848 flags:# gallery:flags.0?true query_id:long next_offset:flags.1?string switch_pm:flags.2?InlineBotSwitchPM results:Vector<BotInlineResult> cache_time:int users:Vector<User> = messages.BotResults;

exportedMessageLink#5dab1af4 link:string html:string = ExportedMessageLink;

messageFwdHeader#5f777dce flags:# from_id:flags.0?Peer from_name:flags.5?string date:int channel_post:flags.2?int post_author:flags.3?string saved_from_peer:flags.4?Peer saved_from_msg_id:flags.4?int psa_type:flags.6?string = MessageFwdHeader;

// @type Type of verification code that will be sent next if you call the resendCode method
// @enum Type of verification code that will be sent next if you call the resendCode method: SMS code
auth.codeTypeSms#72a3158c = auth.CodeType;
// @enum Type of verification code that will be sent next if you call the resendCode method: Phone call
auth.codeTypeCall#741cd3e3 = auth.CodeType;
// @enum Type of verification code that will be sent next if you call the resendCode method: Flash call
auth.codeTypeFlashCall#226ccefb = auth.CodeType;

auth.sentCodeTypeApp#3dbb5986 length:int = auth.SentCodeType;
auth.sentCodeTypeSms#c000bba2 length:int = auth.SentCodeType;
auth.sentCodeTypeCall#5353e5a7 length:int = auth.SentCodeType;
auth.sentCodeTypeFlashCall#ab03c6d9 pattern:string = auth.SentCodeType;

messages.botCallbackAnswer#36585ea4 flags:# alert:flags.1?true has_url:flags.3?true native_ui:flags.4?true message:flags.0?string url:flags.2?string cache_time:int = messages.BotCallbackAnswer;

messages.messageEditData#26b5dde6 flags:# caption:flags.0?true = messages.MessageEditData;

inputBotInlineMessageID#890c3d89 dc_id:int id:long access_hash:long = InputBotInlineMessageID;

inlineBotSwitchPM#3c20629f text:string start_param:string = InlineBotSwitchPM;

messages.peerDialogs#3371c354 dialogs:Vector<Dialog> messages:Vector<Message> chats:Vector<Chat> users:Vector<User> state:updates.State = messages.PeerDialogs;

topPeer#edcdc05b peer:Peer rating:double = TopPeer;

// @type Top peer category
// @enum Most used bots
topPeerCategoryBotsPM#ab661b5b = TopPeerCategory;
// @enum Most used inline bots
topPeerCategoryBotsInline#148677e2 = TopPeerCategory;
// @enum Users we've chatted most frequently with
topPeerCategoryCorrespondents#637b7ed = TopPeerCategory;
// @enum Often-opened groups and supergroups
topPeerCategoryGroups#bd17a14a = TopPeerCategory;
// @enum Most frequently visited channels
topPeerCategoryChannels#161d9628 = TopPeerCategory;
// @enum Most frequently called users
topPeerCategoryPhoneCalls#1e76a78c = TopPeerCategory;
// @enum Users to which the users often forwards messages to
topPeerCategoryForwardUsers#a8406ca9 = TopPeerCategory;
// @enum Chats to which the users often forwards messages to
topPeerCategoryForwardChats#fbeec0f0 = TopPeerCategory;

topPeerCategoryPeers#fb834291 category:TopPeerCategory count:int peers:Vector<TopPeer> = TopPeerCategoryPeers;

contacts.topPeersNotModified#de266ef5 = contacts.TopPeers;
contacts.topPeers#70b772a8 categories:Vector<TopPeerCategoryPeers> chats:Vector<Chat> users:Vector<User> = contacts.TopPeers;
contacts.topPeersDisabled#b52c939d = contacts.TopPeers;

draftMessageEmpty#1b0c841a flags:# date:flags.0?int = DraftMessage;
draftMessage#fd8e711f flags:# no_webpage:flags.1?true reply_to_msg_id:flags.0?int message:string entities:flags.3?Vector<MessageEntity> date:int = DraftMessage;

messages.featuredStickersNotModified#c6dc0c66 count:int = messages.FeaturedStickers;
messages.featuredStickers#b6abc341 hash:int count:int sets:Vector<StickerSetCovered> unread:Vector<long> = messages.FeaturedStickers;

messages.recentStickersNotModified#b17f890 = messages.RecentStickers;
messages.recentStickers#22f3afb3 hash:int packs:Vector<StickerPack> stickers:Vector<Document> dates:Vector<int> = messages.RecentStickers;

messages.archivedStickers#4fcba9c8 count:int sets:Vector<StickerSetCovered> = messages.ArchivedStickers;

messages.stickerSetInstallResultSuccess#38641628 = messages.StickerSetInstallResult;
messages.stickerSetInstallResultArchive#35e410a8 sets:Vector<StickerSetCovered> = messages.StickerSetInstallResult;

stickerSetCovered#6410a5d2 set:StickerSet cover:Document = StickerSetCovered;
stickerSetMultiCovered#3407e51b set:StickerSet covers:Vector<Document> = StickerSetCovered;

maskCoords#aed6dbb2 n:int x:double y:double zoom:double = MaskCoords;

inputStickeredMediaPhoto#4a992157 id:InputPhoto = InputStickeredMedia;
inputStickeredMediaDocument#438865b id:InputDocument = InputStickeredMedia;

game#bdf9653b flags:# id:long access_hash:long short_name:string title:string description:string photo:Photo document:flags.0?Document = Game;

inputGameID#32c3e77 id:long access_hash:long = InputGame;
inputGameShortName#c331e80a bot_id:InputUser short_name:string = InputGame;

highScore#58fffcd0 pos:int user_id:int score:int = HighScore;

messages.highScores#9a3bfd99 scores:Vector<HighScore> users:Vector<User> = messages.HighScores;

textEmpty#dc3d824f = RichText;
textPlain#744694e0 text:string = RichText;
textBold#6724abc4 text:RichText = RichText;
textItalic#d912a59c text:RichText = RichText;
textUnderline#c12622c4 text:RichText = RichText;
textStrike#9bf8bb95 text:RichText = RichText;
textFixed#6c3f19b9 text:RichText = RichText;
textUrl#3c2884c1 text:RichText url:string webpage_id:long = RichText;
textEmail#de5a0dd6 text:RichText email:string = RichText;
textConcat#7e6260d7 texts:Vector<RichText> = RichText;
textSubscript#ed6a8504 text:RichText = RichText;
textSuperscript#c7fb5e01 text:RichText = RichText;
textMarked#34b8621 text:RichText = RichText;
textPhone#1ccb966a text:RichText phone:string = RichText;
textImage#81ccf4f document_id:long w:int h:int = RichText;
textAnchor#35553762 text:RichText name:string = RichText;

pageBlockUnsupported#13567e8a = PageBlock;
pageBlockTitle#70abc3fd text:RichText = PageBlock;
pageBlockSubtitle#8ffa9a1f text:RichText = PageBlock;
pageBlockAuthorDate#baafe5e0 author:RichText published_date:int = PageBlock;
pageBlockHeader#bfd064ec text:RichText = PageBlock;
pageBlockSubheader#f12bb6e1 text:RichText = PageBlock;
pageBlockParagraph#467a0766 text:RichText = PageBlock;
pageBlockPreformatted#c070d93e text:RichText language:string = PageBlock;
pageBlockFooter#48870999 text:RichText = PageBlock;
pageBlockDivider#db20b188 = PageBlock;
pageBlockAnchor#ce0d37b0 name:string = PageBlock;
pageBlockList#e4e88011 items:Vector<PageListItem> = PageBlock;
pageBlockBlockquote#263d7c26 text:RichText caption:RichText = PageBlock;
pageBlockPullquote#4f4456d3 text:RichText caption:RichText = PageBlock;
pageBlockPhoto#1759c560 flags:# photo_id:long caption:PageCaption url:flags.0?string webpage_id:flags.0?long = PageBlock;
pageBlockVideo#7c8fe7b6 flags:# autoplay:flags.0?true loop:flags.1?true video_id:long caption:PageCaption = PageBlock;
pageBlockCover#39f23300 cover:PageBlock = PageBlock;
pageBlockEmbed#a8718dc5 flags:# full_width:flags.0?true allow_scrolling:flags.3?true url:flags.1?string html:flags.2?string poster_photo_id:flags.4?long w:flags.5?int h:flags.5?int caption:PageCaption = PageBlock;
pageBlockEmbedPost#f259a80b url:string webpage_id:long author_photo_id:long author:string date:int blocks:Vector<PageBlock> caption:PageCaption = PageBlock;
pageBlockCollage#65a0fa4d items:Vector<PageBlock> caption:PageCaption = PageBlock;
pageBlockSlideshow#31f9590 items:Vector<PageBlock> caption:PageCaption = PageBlock;
pageBlockChannel#ef1751b5 channel:Chat = PageBlock;
pageBlockAudio#804361ea audio_id:long caption:PageCaption = PageBlock;
pageBlockKicker#1e148390 text:RichText = PageBlock;
pageBlockTable#bf4dea82 flags:# bordered:flags.0?true striped:flags.1?true title:RichText rows:Vector<PageTableRow> = PageBlock;
pageBlockOrderedList#9a8ae1e1 items:Vector<PageListOrderedItem> = PageBlock;
pageBlockDetails#76768bed flags:# open:flags.0?true blocks:Vector<PageBlock> title:RichText = PageBlock;
pageBlockRelatedArticles#16115a96 title:RichText articles:Vector<PageRelatedArticle> = PageBlock;
pageBlockMap#a44f3ef6 geo:GeoPoint zoom:int w:int h:int caption:PageCaption = PageBlock;

// @type Why was the phone call discarded?
// @enum The phone call was missed
phoneCallDiscardReasonMissed#85e42301 = PhoneCallDiscardReason;
// @enum The phone call was disconnected
phoneCallDiscardReasonDisconnect#e095c1a0 = PhoneCallDiscardReason;
// @enum The phone call was ended normally
phoneCallDiscardReasonHangup#57adc690 = PhoneCallDiscardReason;
// @enum The phone call was discared because the user is busy in another call
phoneCallDiscardReasonBusy#faf7e8c9 = PhoneCallDiscardReason;

dataJSON#7d748d04 data:string = DataJSON;

labeledPrice#cb296bf8 label:string amount:long = LabeledPrice;

invoice#c30aa358 flags:# test:flags.0?true name_requested:flags.1?true phone_requested:flags.2?true email_requested:flags.3?true shipping_address_requested:flags.4?true flexible:flags.5?true phone_to_provider:flags.6?true email_to_provider:flags.7?true currency:string prices:Vector<LabeledPrice> = Invoice;

paymentCharge#ea02c27e id:string provider_charge_id:string = PaymentCharge;

postAddress#1e8caaeb street_line1:string street_line2:string city:string state:string country_iso2:string post_code:string = PostAddress;

paymentRequestedInfo#909c3f94 flags:# name:flags.0?string phone:flags.1?string email:flags.2?string shipping_address:flags.3?PostAddress = PaymentRequestedInfo;

paymentSavedCredentialsCard#cdc27a1f id:string title:string = PaymentSavedCredentials;

webDocument#1c570ed1 url:string access_hash:long size:int mime_type:string attributes:Vector<DocumentAttribute> = WebDocument;
webDocumentNoProxy#f9c8bcc6 url:string size:int mime_type:string attributes:Vector<DocumentAttribute> = WebDocument;

inputWebDocument#9bed434d url:string size:int mime_type:string attributes:Vector<DocumentAttribute> = InputWebDocument;

inputWebFileLocation#c239d686 url:string access_hash:long = InputWebFileLocation;
inputWebFileGeoPointLocation#9f2221c9 geo_point:InputGeoPoint access_hash:long w:int h:int zoom:int scale:int = InputWebFileLocation;

upload.webFile#21e753bc size:int mime_type:string file_type:storage.FileType mtime:int bytes:bytes = upload.WebFile;

// @constructor Payment form
// @param bot_id               Bot ID
// @param can_save_credentials Whether the user can choose to save credentials.
// @param flags                Flags, see TL conditional fields
// @param invoice              Invoice
// @param native_params        Contains information about the payment provider, if available, to support it natively without the need for opening the URL.A JSON object that can contain the following fields:- publishable_key: Stripe API publishable key- apple_pay_merchant_id: Apple Pay merchant ID- android_pay_public_key: Android Pay public key- android_pay_bgcolor: Android Pay form background color- android_pay_inverse: Whether to use the dark theme in the Android Pay form- need_country: True, if the user country must be provided,- need_zip: True, if the user ZIP/postal code must be provided,- need_cardholder_name: True, if the cardholder name must be provided
// @param native_provider      Payment provider name.One of the following:- stripe
// @param password_missing     Indicates that the user can save payment credentials, but only after setting up a 2FA password (currently the account doesn't have a 2FA password)
// @param provider_id          Payment provider ID.
// @param saved_credentials    Contains information about saved card credentials
// @param saved_info           Saved server-side order information
// @param url                  Payment form URL
// @param users                Users
payments.paymentForm#3f56aea3 flags:# can_save_credentials:flags.2?true password_missing:flags.3?true bot_id:int invoice:Invoice provider_id:int url:string native_provider:flags.4?string native_params:flags.4?DataJSON saved_info:flags.0?PaymentRequestedInfo saved_credentials:flags.1?PaymentSavedCredentials users:Vector<User> = payments.PaymentForm;

payments.validatedRequestedInfo#d1451883 flags:# id:flags.0?string shipping_options:flags.1?Vector<ShippingOption> = payments.ValidatedRequestedInfo;

payments.paymentResult#4e5f810d updates:Updates = payments.PaymentResult;
payments.paymentVerificationNeeded#d8411139 url:string = payments.PaymentResult;

payments.paymentReceipt#500911e1 flags:# date:int bot_id:int invoice:Invoice provider_id:int info:flags.0?PaymentRequestedInfo shipping:flags.1?ShippingOption currency:string total_amount:long credentials_title:string users:Vector<User> = payments.PaymentReceipt;

payments.savedInfo#fb8fe43c flags:# has_saved_credentials:flags.1?true saved_info:flags.0?PaymentRequestedInfo = payments.SavedInfo;

inputPaymentCredentialsSaved#c10eb2cf id:string tmp_password:bytes = InputPaymentCredentials;
inputPaymentCredentials#3417d728 flags:# save:flags.0?true data:DataJSON = InputPaymentCredentials;
inputPaymentCredentialsApplePay#aa1c39f payment_data:DataJSON = InputPaymentCredentials;
inputPaymentCredentialsAndroidPay#ca05d50e payment_token:DataJSON google_transaction_id:string = InputPaymentCredentials;

account.tmpPassword#db64fd34 tmp_password:bytes valid_until:int = account.TmpPassword;

shippingOption#b6213cdf id:string title:string prices:Vector<LabeledPrice> = ShippingOption;

inputStickerSetItem#ffa0a496 flags:# document:InputDocument emoji:string mask_coords:flags.0?MaskCoords = InputStickerSetItem;

inputPhoneCall#1e36fded id:long access_hash:long = InputPhoneCall;

phoneCallEmpty#5366c915 id:long = PhoneCall;
phoneCallWaiting#1b8f4ad1 flags:# video:flags.6?true id:long access_hash:long date:int admin_id:int participant_id:int protocol:PhoneCallProtocol receive_date:flags.0?int = PhoneCall;
phoneCallRequested#87eabb53 flags:# video:flags.6?true id:long access_hash:long date:int admin_id:int participant_id:int g_a_hash:bytes protocol:PhoneCallProtocol = PhoneCall;
phoneCallAccepted#997c454a flags:# video:flags.6?true id:long access_hash:long date:int admin_id:int participant_id:int g_b:bytes protocol:PhoneCallProtocol = PhoneCall;
phoneCall#8742ae7f flags:# p2p_allowed:flags.5?true video:flags.6?true id:long access_hash:long date:int admin_id:int participant_id:int g_a_or_b:bytes key_fingerprint:long protocol:PhoneCallProtocol connections:Vector<PhoneConnection> start_date:int = PhoneCall;
phoneCallDiscarded#50ca4de1 flags:# need_rating:flags.2?true need_debug:flags.3?true video:flags.6?true id:long reason:flags.0?PhoneCallDiscardReason duration:flags.1?int = PhoneCall;

phoneConnection#9d4c17c0 id:long ip:string ipv6:string port:int peer_tag:bytes = PhoneConnection;
phoneConnectionWebrtc#635fe375 flags:# turn:flags.0?true stun:flags.1?true id:long ip:string ipv6:string port:int username:string password:string = PhoneConnection;

phoneCallProtocol#fc878fc8 flags:# udp_p2p:flags.0?true udp_reflector:flags.1?true min_layer:int max_layer:int library_versions:Vector<string> = PhoneCallProtocol;

phone.phoneCall#ec82e140 phone_call:PhoneCall users:Vector<User> = phone.PhoneCall;

upload.cdnFileReuploadNeeded#eea8e46e request_token:bytes = upload.CdnFile;
upload.cdnFile#a99fca4f bytes:bytes = upload.CdnFile;

cdnPublicKey#c982eaba dc_id:int public_key:string = CdnPublicKey;

cdnConfig#5725e40a public_keys:Vector<CdnPublicKey> = CdnConfig;

langPackString#cad181f6 key:string value:string = LangPackString;
langPackStringPluralized#6c47ac9f flags:# key:string zero_value:flags.0?string one_value:flags.1?string two_value:flags.2?string few_value:flags.3?string many_value:flags.4?string other_value:string = LangPackString;
langPackStringDeleted#2979eeb2 key:string = LangPackString;

langPackDifference#f385c1f6 lang_code:string from_version:int version:int strings:Vector<LangPackString> = LangPackDifference;

langPackLanguage#eeca5ce3 flags:# official:flags.0?true rtl:flags.2?true beta:flags.3?true name:string native_name:string lang_code:string base_lang_code:flags.1?string plural_code:string strings_count:int translated_count:int translations_url:string = LangPackLanguage;

channelAdminLogEventActionChangeTitle#e6dfb825 prev_value:string new_value:string = ChannelAdminLogEventAction;
channelAdminLogEventActionChangeAbout#55188a2e prev_value:string new_value:string = ChannelAdminLogEventAction;
channelAdminLogEventActionChangeUsername#6a4afc38 prev_value:string new_value:string = ChannelAdminLogEventAction;
channelAdminLogEventActionChangePhoto#434bd2af prev_photo:Photo new_photo:Photo = ChannelAdminLogEventAction;
channelAdminLogEventActionToggleInvites#1b7907ae new_value:Bool = ChannelAdminLogEventAction;
channelAdminLogEventActionToggleSignatures#26ae0971 new_value:Bool = ChannelAdminLogEventAction;
channelAdminLogEventActionUpdatePinned#e9e82c18 message:Message = ChannelAdminLogEventAction;
channelAdminLogEventActionEditMessage#709b2405 prev_message:Message new_message:Message = ChannelAdminLogEventAction;
channelAdminLogEventActionDeleteMessage#42e047bb message:Message = ChannelAdminLogEventAction;
channelAdminLogEventActionParticipantJoin#183040d3 = ChannelAdminLogEventAction;
channelAdminLogEventActionParticipantLeave#f89777f2 = ChannelAdminLogEventAction;
channelAdminLogEventActionParticipantInvite#e31c34d8 participant:ChannelParticipant = ChannelAdminLogEventAction;
channelAdminLogEventActionParticipantToggleBan#e6d83d7e prev_participant:ChannelParticipant new_participant:ChannelParticipant = ChannelAdminLogEventAction;
channelAdminLogEventActionParticipantToggleAdmin#d5676710 prev_participant:ChannelParticipant new_participant:ChannelParticipant = ChannelAdminLogEventAction;
channelAdminLogEventActionChangeStickerSet#b1c3caa7 prev_stickerset:InputStickerSet new_stickerset:InputStickerSet = ChannelAdminLogEventAction;
channelAdminLogEventActionTogglePreHistoryHidden#5f5c95f1 new_value:Bool = ChannelAdminLogEventAction;
channelAdminLogEventActionDefaultBannedRights#2df5fc0a prev_banned_rights:ChatBannedRights new_banned_rights:ChatBannedRights = ChannelAdminLogEventAction;
channelAdminLogEventActionStopPoll#8f079643 message:Message = ChannelAdminLogEventAction;
channelAdminLogEventActionChangeLinkedChat#a26f881b prev_value:int new_value:int = ChannelAdminLogEventAction;
channelAdminLogEventActionChangeLocation#e6b76ae prev_value:ChannelLocation new_value:ChannelLocation = ChannelAdminLogEventAction;
channelAdminLogEventActionToggleSlowMode#53909779 prev_value:int new_value:int = ChannelAdminLogEventAction;

channelAdminLogEvent#3b5a3e40 id:long date:int user_id:int action:ChannelAdminLogEventAction = ChannelAdminLogEvent;

channels.adminLogResults#ed8af74d events:Vector<ChannelAdminLogEvent> chats:Vector<Chat> users:Vector<User> = channels.AdminLogResults;

channelAdminLogEventsFilter#ea107ae4 flags:# join:flags.0?true leave:flags.1?true invite:flags.2?true ban:flags.3?true unban:flags.4?true kick:flags.5?true unkick:flags.6?true promote:flags.7?true demote:flags.8?true info:flags.9?true settings:flags.10?true pinned:flags.11?true edit:flags.12?true delete:flags.13?true = ChannelAdminLogEventsFilter;

popularContact#5ce14175 client_id:long importers:int = PopularContact;

messages.favedStickersNotModified#9e8fa6d3 = messages.FavedStickers;
messages.favedStickers#f37f2f16 hash:int packs:Vector<StickerPack> stickers:Vector<Document> = messages.FavedStickers;

recentMeUrlUnknown#46e1d13d url:string = RecentMeUrl;
recentMeUrlUser#8dbc3336 url:string user_id:int = RecentMeUrl;
recentMeUrlChat#a01b22f9 url:string chat_id:int = RecentMeUrl;
recentMeUrlChatInvite#eb49081d url:string chat_invite:ChatInvite = RecentMeUrl;
recentMeUrlStickerSet#bc0a57dc url:string set:StickerSetCovered = RecentMeUrl;

help.recentMeUrls#e0310d7 urls:Vector<RecentMeUrl> chats:Vector<Chat> users:Vector<User> = help.RecentMeUrls;

inputSingleMedia#1cc6e91f flags:# media:InputMedia random_id:long message:string entities:flags.0?Vector<MessageEntity> = InputSingleMedia;

webAuthorization#cac943f2 hash:long bot_id:int domain:string browser:string platform:string date_created:int date_active:int ip:string region:string = WebAuthorization;

account.webAuthorizations#ed56c9fc authorizations:Vector<WebAuthorization> users:Vector<User> = account.WebAuthorizations;

inputMessageID#a676a322 id:int = InputMessage;
inputMessageReplyTo#bad88395 id:int = InputMessage;
inputMessagePinned#86872538 = InputMessage;
inputMessageCallbackQuery#acfa1a7e id:int query_id:long = InputMessage;

inputDialogPeer#fcaafeb7 peer:InputPeer = InputDialogPeer;
inputDialogPeerFolder#64600527 folder_id:int = InputDialogPeer;

dialogPeer#e56dbf05 peer:Peer = DialogPeer;
dialogPeerFolder#514519e2 folder_id:int = DialogPeer;

messages.foundStickerSetsNotModified#d54b65d = messages.FoundStickerSets;
messages.foundStickerSets#5108d648 hash:int sets:Vector<StickerSetCovered> = messages.FoundStickerSets;

fileHash#6242c773 offset:int limit:int hash:bytes = FileHash;

inputClientProxy#75588b3f address:string port:int = InputClientProxy;

help.termsOfServiceUpdateEmpty#e3309f7f expires:int = help.TermsOfServiceUpdate;
help.termsOfServiceUpdate#28ecf961 expires:int terms_of_service:help.TermsOfService = help.TermsOfServiceUpdate;

inputSecureFileUploaded#3334b0f0 id:long parts:int md5_checksum:string file_hash:bytes secret:bytes = InputSecureFile;
inputSecureFile#5367e5be id:long access_hash:long = InputSecureFile;

secureFileEmpty#64199744 = SecureFile;
secureFile#e0277a62 id:long access_hash:long size:int dc_id:int date:int file_hash:bytes secret:bytes = SecureFile;

secureData#8aeabec3 data:bytes data_hash:bytes secret:bytes = SecureData;

securePlainPhone#7d6099dd phone:string = SecurePlainData;
securePlainEmail#21ec5a5f email:string = SecurePlainData;

// @type Secure value type
// @enum Personal details
secureValueTypePersonalDetails#9d2a81e3 = SecureValueType;
// @enum Passport
secureValueTypePassport#3dac6a00 = SecureValueType;
// @enum Driver's license
secureValueTypeDriverLicense#6e425c4 = SecureValueType;
// @enum Identity card
secureValueTypeIdentityCard#a0d0744b = SecureValueType;
// @enum Internal [passport](https://core.telegram.org/passport)
secureValueTypeInternalPassport#99a48f23 = SecureValueType;
// @enum Address
secureValueTypeAddress#cbe31e26 = SecureValueType;
// @enum Utility bill
secureValueTypeUtilityBill#fc36954e = SecureValueType;
// @enum Bank statement
secureValueTypeBankStatement#89137c0d = SecureValueType;
// @enum Rental agreement
secureValueTypeRentalAgreement#8b883488 = SecureValueType;
// @enum Internal registration [passport](https://core.telegram.org/passport)
secureValueTypePassportRegistration#99e3806a = SecureValueType;
// @enum Temporary registration
secureValueTypeTemporaryRegistration#ea02ec33 = SecureValueType;
// @enum Phone
secureValueTypePhone#b320aadb = SecureValueType;
// @enum Email
secureValueTypeEmail#8e3ca7ee = SecureValueType;

secureValue#187fa0ca flags:# type:SecureValueType data:flags.0?SecureData front_side:flags.1?SecureFile reverse_side:flags.2?SecureFile selfie:flags.3?SecureFile translation:flags.6?Vector<SecureFile> files:flags.4?Vector<SecureFile> plain_data:flags.5?SecurePlainData hash:bytes = SecureValue;

inputSecureValue#db21d0a7 flags:# type:SecureValueType data:flags.0?SecureData front_side:flags.1?InputSecureFile reverse_side:flags.2?InputSecureFile selfie:flags.3?InputSecureFile translation:flags.6?Vector<InputSecureFile> files:flags.4?Vector<InputSecureFile> plain_data:flags.5?SecurePlainData = InputSecureValue;

secureValueHash#ed1ecdb0 type:SecureValueType hash:bytes = SecureValueHash;

secureValueErrorData#e8a40bd9 type:SecureValueType data_hash:bytes field:string text:string = SecureValueError;
secureValueErrorFrontSide#be3dfa type:SecureValueType file_hash:bytes text:string = SecureValueError;
secureValueErrorReverseSide#868a2aa5 type:SecureValueType file_hash:bytes text:string = SecureValueError;
secureValueErrorSelfie#e537ced6 type:SecureValueType file_hash:bytes text:string = SecureValueError;
secureValueErrorFile#7a700873 type:SecureValueType file_hash:bytes text:string = SecureValueError;
secureValueErrorFiles#666220e9 type:SecureValueType file_hash:Vector<bytes> text:string = SecureValueError;
secureValueError#869d758f type:SecureValueType hash:bytes text:string = SecureValueError;
secureValueErrorTranslationFile#a1144770 type:SecureValueType file_hash:bytes text:string = SecureValueError;
secureValueErrorTranslationFiles#34636dd8 type:SecureValueType file_hash:Vector<bytes> text:string = SecureValueError;

// @type Encrypted secure credentials
// @constructor Encrypted credentials required to decrypt telegram passport data.
secureCredentialsEncrypted#33f0ea47 data:bytes hash:bytes secret:bytes = SecureCredentialsEncrypted;

// @type Authorization form
// @constructor Telegram Passport authorization form
// @param privacy_policy_url URL of the service's privacy policy
// @param required_types     Required Telegram Passport¹ documents
// @param users              Info about the bot to which the form will be submitted
// @param values             Already submitted Telegram Passport¹ documents
account.authorizationForm#ad2e1cd8 flags:# required_types:Vector<SecureRequiredType> values:Vector<SecureValue> errors:Vector<SecureValueError> users:Vector<User> privacy_policy_url:flags.0?string = account.AuthorizationForm;

account.sentEmailCode#811f854f email_pattern:string length:int = account.SentEmailCode;

help.deepLinkInfoEmpty#66afa166 = help.DeepLinkInfo;
help.deepLinkInfo#6a4ee832 flags:# update_app:flags.0?true message:string entities:flags.1?Vector<MessageEntity> = help.DeepLinkInfo;

savedPhoneContact#1142bd56 phone:string first_name:string last_name:string date:int = SavedContact;

account.takeout#4dba4501 id:long = account.Takeout;

passwordKdfAlgoUnknown#d45ab096 = PasswordKdfAlgo;
passwordKdfAlgoSHA256SHA256PBKDF2HMACSHA512iter100000SHA256ModPow#3a912d4a salt1:bytes salt2:bytes g:int p:bytes = PasswordKdfAlgo;

securePasswordKdfAlgoUnknown#4a8537 = SecurePasswordKdfAlgo;
securePasswordKdfAlgoPBKDF2HMACSHA512iter100000#bbf2dda0 salt:bytes = SecurePasswordKdfAlgo;
securePasswordKdfAlgoSHA512#86471d92 salt:bytes = SecurePasswordKdfAlgo;

secureSecretSettings#1527bcac secure_algo:SecurePasswordKdfAlgo secure_secret:bytes secure_secret_id:long = SecureSecretSettings;

inputCheckPasswordEmpty#9880f658 = InputCheckPasswordSRP;
inputCheckPasswordSRP#d27ff082 srp_id:long A:bytes M1:bytes = InputCheckPasswordSRP;

secureRequiredType#829d99da flags:# native_names:flags.0?true selfie_required:flags.1?true translation_required:flags.2?true type:SecureValueType = SecureRequiredType;
secureRequiredTypeOneOf#27477b4 types:Vector<SecureRequiredType> = SecureRequiredType;

help.passportConfigNotModified#bfb9f457 = help.PassportConfig;
help.passportConfig#a098d6af hash:int countries_langs:DataJSON = help.PassportConfig;

inputAppEvent#1d1b1245 time:double type:string peer:long data:JSONValue = InputAppEvent;

jsonObjectValue#c0de1bd9 key:string value:JSONValue = JSONObjectValue;

jsonNull#3f6d7b68 = JSONValue;
jsonBool#c7345e6a value:Bool = JSONValue;
jsonNumber#2be0dfa4 value:double = JSONValue;
jsonString#b71e767a value:string = JSONValue;
jsonArray#f7444763 value:Vector<JSONValue> = JSONValue;
jsonObject#99c1d49d value:Vector<JSONObjectValue> = JSONValue;

pageTableCell#34566b6a flags:# header:flags.0?true align_center:flags.3?true align_right:flags.4?true valign_middle:flags.5?true valign_bottom:flags.6?true text:flags.7?RichText colspan:flags.1?int rowspan:flags.2?int = PageTableCell;

pageTableRow#e0c0c5e5 cells:Vector<PageTableCell> = PageTableRow;

pageCaption#6f747657 text:RichText credit:RichText = PageCaption;

pageListItemText#b92fb6cd text:RichText = PageListItem;
pageListItemBlocks#25e073fc blocks:Vector<PageBlock> = PageListItem;

pageListOrderedItemText#5e068047 num:string text:RichText = PageListOrderedItem;
pageListOrderedItemBlocks#98dd8936 num:string blocks:Vector<PageBlock> = PageListOrderedItem;

pageRelatedArticle#b390dc08 flags:# url:string webpage_id:long title:flags.0?string description:flags.1?string photo_id:flags.2?long author:flags.3?string published_date:flags.4?int = PageRelatedArticle;

page#98657f0d flags:# part:flags.0?true rtl:flags.1?true v2:flags.2?true url:string blocks:Vector<PageBlock> photos:Vector<Photo> documents:Vector<Document> views:flags.3?int = Page;

help.supportName#8c05f1c9 name:string = help.SupportName;

help.userInfoEmpty#f3ae2eed = help.UserInfo;
help.userInfo#1eb3758 message:string entities:Vector<MessageEntity> author:string date:int = help.UserInfo;

pollAnswer#6ca9c2e9 text:string option:bytes = PollAnswer;

poll#86e18161 id:long flags:# closed:flags.0?true public_voters:flags.1?true multiple_choice:flags.2?true quiz:flags.3?true question:string answers:Vector<PollAnswer> close_period:flags.4?int close_date:flags.5?int = Poll;

pollAnswerVoters#3b6ddad2 flags:# chosen:flags.0?true correct:flags.1?true option:bytes voters:int = PollAnswerVoters;

pollResults#badcc1a3 flags:# min:flags.0?true results:flags.1?Vector<PollAnswerVoters> total_voters:flags.2?int recent_voters:flags.3?Vector<int> solution:flags.4?string solution_entities:flags.4?Vector<MessageEntity> = PollResults;

chatOnlines#f041e250 onlines:int = ChatOnlines;

statsURL#47a971e0 url:string = StatsURL;

chatAdminRights#5fb224d5 flags:# change_info:flags.0?true post_messages:flags.1?true edit_messages:flags.2?true delete_messages:flags.3?true ban_users:flags.4?true invite_users:flags.5?true pin_messages:flags.7?true add_admins:flags.9?true anonymous:flags.10?true manage_call:flags.11?true = ChatAdminRights;

chatBannedRights#9f120418 flags:# view_messages:flags.0?true send_messages:flags.1?true send_media:flags.2?true send_stickers:flags.3?true send_gifs:flags.4?true send_games:flags.5?true send_inline:flags.6?true embed_links:flags.7?true send_polls:flags.8?true change_info:flags.10?true invite_users:flags.15?true pin_messages:flags.17?true until_date:int = ChatBannedRights;

inputWallPaper#e630b979 id:long access_hash:long = InputWallPaper;
inputWallPaperSlug#72091c80 slug:string = InputWallPaper;
inputWallPaperNoFile#8427bbac = InputWallPaper;

account.wallPapersNotModified#1c199183 = account.WallPapers;
account.wallPapers#702b65a9 hash:int wallpapers:Vector<WallPaper> = account.WallPapers;

codeSettings#debebe83 flags:# allow_flashcall:flags.0?true current_number:flags.1?true allow_app_hash:flags.4?true = CodeSettings;

wallPaperSettings#5086cf8 flags:# blur:flags.1?true motion:flags.2?true background_color:flags.0?int second_background_color:flags.4?int intensity:flags.3?int rotation:flags.4?int = WallPaperSettings;

autoDownloadSettings#e04232f3 flags:# disabled:flags.0?true video_preload_large:flags.1?true audio_preload_next:flags.2?true phonecalls_less_data:flags.3?true photo_size_max:int video_size_max:int file_size_max:int video_upload_maxbitrate:int = AutoDownloadSettings;

account.autoDownloadSettings#63cacf26 low:AutoDownloadSettings medium:AutoDownloadSettings high:AutoDownloadSettings = account.AutoDownloadSettings;

emojiKeyword#d5b3b9f9 keyword:string emoticons:Vector<string> = EmojiKeyword;
emojiKeywordDeleted#236df622 keyword:string emoticons:Vector<string> = EmojiKeyword;

emojiKeywordsDifference#5cc761bd lang_code:string from_version:int version:int keywords:Vector<EmojiKeyword> = EmojiKeywordsDifference;

emojiURL#a575739d url:string = EmojiURL;

emojiLanguage#b3fb5361 lang_code:string = EmojiLanguage;

fileLocationToBeDeprecated#bc7fc6cd volume_id:long local_id:int = FileLocation;

folder#ff544e65 flags:# autofill_new_broadcasts:flags.0?true autofill_public_groups:flags.1?true autofill_new_correspondents:flags.2?true id:int title:string photo:flags.3?ChatPhoto = Folder;

inputFolderPeer#fbd2c296 peer:InputPeer folder_id:int = InputFolderPeer;

folderPeer#e9baa668 peer:Peer folder_id:int = FolderPeer;

messages.searchCounter#e844ebff flags:# inexact:flags.1?true filter:MessagesFilter count:int = messages.SearchCounter;

urlAuthResultRequest#92d33a0e flags:# request_write_access:flags.0?true bot:User domain:string = UrlAuthResult;
urlAuthResultAccepted#8f8c0e4e url:string = UrlAuthResult;
urlAuthResultDefault#a9d6db1f = UrlAuthResult;

channelLocationEmpty#bfb5ad8b = ChannelLocation;
channelLocation#209b82db geo_point:GeoPoint address:string = ChannelLocation;

peerLocated#ca461b5d peer:Peer expires:int distance:int = PeerLocated;
peerSelfLocated#f8ec284b expires:int = PeerLocated;

restrictionReason#d072acb4 platform:string reason:string text:string = RestrictionReason;

inputTheme#3c5693e9 id:long access_hash:long = InputTheme;
inputThemeSlug#f5890df1 slug:string = InputTheme;

theme#28f1114 flags:# creator:flags.0?true default:flags.1?true id:long access_hash:long slug:string title:string document:flags.2?Document settings:flags.3?ThemeSettings installs_count:int = Theme;

account.themesNotModified#f41eb622 = account.Themes;
account.themes#7f676421 hash:int themes:Vector<Theme> = account.Themes;

auth.loginToken#629f1980 expires:int token:bytes = auth.LoginToken;
auth.loginTokenMigrateTo#68e9916 dc_id:int token:bytes = auth.LoginToken;
auth.loginTokenSuccess#390d5c5e authorization:auth.Authorization = auth.LoginToken;

account.contentSettings#57e28221 flags:# sensitive_enabled:flags.0?true sensitive_can_change:flags.1?true = account.ContentSettings;

messages.inactiveChats#a927fec5 dates:Vector<int> chats:Vector<Chat> users:Vector<User> = messages.InactiveChats;

// @type Basic theme settings
// @enum Classic theme
baseThemeClassic#c3a12462 = BaseTheme;
// @enum Day theme
baseThemeDay#fbd81688 = BaseTheme;
// @enum Night theme
baseThemeNight#b7b31ea8 = BaseTheme;
// @enum Tinted theme
baseThemeTinted#6d5f77ee = BaseTheme;
// @enum Arctic theme
baseThemeArctic#5b11125a = BaseTheme;

inputThemeSettings#bd507cd1 flags:# base_theme:BaseTheme accent_color:int message_top_color:flags.0?int message_bottom_color:flags.0?int wallpaper:flags.1?InputWallPaper wallpaper_settings:flags.1?WallPaperSettings = InputThemeSettings;

themeSettings#9c14984a flags:# base_theme:BaseTheme accent_color:int message_top_color:flags.0?int message_bottom_color:flags.0?int wallpaper:flags.1?WallPaper = ThemeSettings;

webPageAttributeTheme#54b56617 flags:# documents:flags.0?Vector<Document> settings:flags.1?ThemeSettings = WebPageAttribute;

messageUserVote#a28e5559 user_id:int option:bytes date:int = MessageUserVote;
messageUserVoteInputOption#36377430 user_id:int date:int = MessageUserVote;
messageUserVoteMultiple#e8fe0de user_id:int options:Vector<bytes> date:int = MessageUserVote;

messages.votesList#823f649 flags:# count:int votes:Vector<MessageUserVote> users:Vector<User> next_offset:flags.0?string = messages.VotesList;

bankCardOpenUrl#f568028a url:string name:string = BankCardOpenUrl;

payments.bankCardData#3e24e573 title:string open_urls:Vector<BankCardOpenUrl> = payments.BankCardData;

dialogFilter#7438f7e8 flags:# contacts:flags.0?true non_contacts:flags.1?true groups:flags.2?true broadcasts:flags.3?true bots:flags.4?true exclude_muted:flags.11?true exclude_read:flags.12?true exclude_archived:flags.13?true id:int title:string emoticon:flags.25?string pinned_peers:Vector<InputPeer> include_peers:Vector<InputPeer> exclude_peers:Vector<InputPeer> = DialogFilter;

dialogFilterSuggested#77744d4a filter:DialogFilter description:string = DialogFilterSuggested;

statsDateRangeDays#b637edaf min_date:int max_date:int = StatsDateRangeDays;

statsAbsValueAndPrev#cb43acde current:double previous:double = StatsAbsValueAndPrev;

statsPercentValue#cbce2fe0 part:double total:double = StatsPercentValue;

// @type Channel statistics graph
// @constructor This [channel statistics graph](https://core.telegram.org/api/stats) must be generated asynchronously using StatsLoadAsyncGraph to reduce server load
// @param token Token to use for fetching the async graph
statsGraphAsync#4a27eb2d token:string = StatsGraph;
// @constructor An error occurred while generating the statistics graph
// @param error The error
statsGraphError#bedc9822 error:string = StatsGraph;
// @constructor
// @param json       Statistics data
// @param zoom_token Zoom token
statsGraph#8ea464b6 flags:# json:DataJSON zoom_token:flags.0?string = StatsGraph;

messageInteractionCounters#ad4fc9bd msg_id:int views:int forwards:int = MessageInteractionCounters;

// @type Channel statistics
// @constructor [Channel statistics](https://core.telegram.org/api/stats)
// @param period                        Period in consideration
// @param followers                     Follower count change for period in consideration
// @param views_per_post                `total_viewcount/postcount`, for posts posted during the period in consideration (`views_per_post`). Note that in this case, `current` refers to the `period` in consideration (`min_date` till `max_date`), and `prev` refers to the previous period (`(min_date - (max_date - min_date))` till `min_date`)
// @param shares_per_post               `total_viewcount/postcount`, for posts posted during the period in consideration (`views_per_post`). Note that in this case, `current` refers to the `period` in consideration (`min_date` till `max_date`), and `prev` refers to the previous period (`(min_date - (max_date - min_date))` till `min_date`)
// @param enabled_notifications         Percentage of subscribers with enabled notifications
// @param growth_graph                  Channel growth graph (absolute subscriber count)
// @param followers_graph               Followers growth graph (relative subscriber count)
// @param mute_graph                    Muted users graph (relative)
// @param top_hours_graph               Views per hour graph (absolute)
// @param interactions_graph            Interactions graph (absolute)
// @param iv_interactions_graph         IV interactions graph (absolute)
// @param views_by_source_graph         Views by source graph (absolute)
// @param new_followers_by_source_graph New followers by source graph (absolute)
// @param languages_graph               Subscriber language graph (piechart)
// @param recent_message_interactions   Recent message interactions
stats.broadcastStats#bdf78394 period:StatsDateRangeDays followers:StatsAbsValueAndPrev views_per_post:StatsAbsValueAndPrev shares_per_post:StatsAbsValueAndPrev enabled_notifications:StatsPercentValue growth_graph:StatsGraph followers_graph:StatsGraph mute_graph:StatsGraph top_hours_graph:StatsGraph interactions_graph:StatsGraph iv_interactions_graph:StatsGraph views_by_source_graph:StatsGraph new_followers_by_source_graph:StatsGraph languages_graph:StatsGraph recent_message_interactions:Vector<MessageInteractionCounters> = stats.BroadcastStats;

help.promoDataEmpty#98f6ac75 expires:int = help.PromoData;
help.promoData#8c39793f flags:# proxy:flags.0?true expires:int peer:Peer chats:Vector<Chat> users:Vector<User> psa_type:flags.1?string psa_message:flags.2?string = help.PromoData;

videoSize#e831c556 flags:# type:string location:FileLocation w:int h:int size:int video_start_ts:flags.0?double = VideoSize;

statsGroupTopPoster#18f3d0f7 user_id:int messages:int avg_chars:int = StatsGroupTopPoster;

statsGroupTopAdmin#6014f412 user_id:int deleted:int kicked:int banned:int = StatsGroupTopAdmin;

statsGroupTopInviter#31962a4c user_id:int invitations:int = StatsGroupTopInviter;

stats.megagroupStats#ef7ff916 period:StatsDateRangeDays members:StatsAbsValueAndPrev messages:StatsAbsValueAndPrev viewers:StatsAbsValueAndPrev posters:StatsAbsValueAndPrev growth_graph:StatsGraph members_graph:StatsGraph new_members_by_source_graph:StatsGraph languages_graph:StatsGraph messages_graph:StatsGraph actions_graph:StatsGraph top_hours_graph:StatsGraph weekdays_graph:StatsGraph top_posters:Vector<StatsGroupTopPoster> top_admins:Vector<StatsGroupTopAdmin> top_inviters:Vector<StatsGroupTopInviter> users:Vector<User> = stats.MegagroupStats;

globalPrivacySettings#bea2f424 flags:# archive_and_mute_new_noncontact_peers:flags.0?Bool = GlobalPrivacySettings;

help.countryCode#4203c5ef flags:# country_code:string prefixes:flags.0?Vector<string> patterns:flags.1?Vector<string> = help.CountryCode;

help.country#c3878e23 flags:# hidden:flags.0?true iso2:string default_name:string name:flags.1?string country_codes:Vector<help.CountryCode> = help.Country;

help.countriesListNotModified#93cc1f32 = help.CountriesList;
help.countriesList#87d0759e countries:Vector<help.Country> hash:int = help.CountriesList;

messageViews#455b853d flags:# views:flags.0?int forwards:flags.1?int replies:flags.2?MessageReplies = MessageViews;

messages.messageViews#b6c4f543 views:Vector<MessageViews> chats:Vector<Chat> users:Vector<User> = messages.MessageViews;

messages.discussionMessage#f5dd8f9d flags:# messages:Vector<Message> max_id:flags.0?int read_inbox_max_id:flags.1?int read_outbox_max_id:flags.2?int chats:Vector<Chat> users:Vector<User> = messages.DiscussionMessage;

messageReplyHeader#a6d57763 flags:# reply_to_msg_id:int reply_to_peer_id:flags.0?Peer reply_to_top_id:flags.1?int = MessageReplyHeader;

messageReplies#4128faac flags:# comments:flags.0?true replies:int replies_pts:int recent_repliers:flags.1?Vector<Peer> channel_id:flags.0?int max_id:flags.2?int read_max_id:flags.3?int = MessageReplies;

peerBlocked#e8fd8014 peer_id:Peer date:int = PeerBlocked;

stats.messageStats#8999f295 views_graph:StatsGraph = stats.MessageStats;

---functions---

// @method Invokes a query after successfull completion of one of the previous queries.
// @param msg_id Message identifier on which a current query depends
// @param query The query itself
invokeAfterMsg#cb9f372d {X:Type} msg_id:long query:!X = X;
// @method Invokes a query after a successfull completion of previous queries
// @param msg_ids List of messages on which a current query depends
// @param query The query itself
invokeAfterMsgs#3dc4b4f0 {X:Type} msg_ids:Vector<long> query:!X = X;
// @method Initialize connection
// @param flags Flags, see TL conditional fields
// @param api_id Application identifier (see. App configuration)
// @param device_model Device model
// @param system_version Operation system version
// @param app_version Application version
// @param system_lang_code Code for the language used on the device's OS, ISO 639-1 standard
// @param lang_pack Language pack to use
// @param lang_code Code for the language used on the client, ISO 639-1 standard
// @param proxy Info about an MTProto proxy
// @param params Additional initConnection parameters. For now, only the tz_offset field is supported, for specifying timezone offset in seconds.
// @param query The query itself
initConnection#c1cd5ea9 {X:Type} flags:# api_id:int device_model:string system_version:string app_version:string system_lang_code:string lang_pack:string lang_code:string proxy:flags.0?InputClientProxy params:flags.1?JSONValue query:!X = X;
// @method Invoke the specified query using the specified API layer
// @param layer The layer to use
// @param query The query
invokeWithLayer#da9b0d0d {X:Type} layer:int query:!X = X;
// @method Invoke a request without subscribing the used connection for updates (this is enabled by default for file queries).
// @param query The query
invokeWithoutUpdates#bf9459b7 {X:Type} query:!X = X;
// @method Invoke with the given message range
// @param range Message range
// @param query Query
invokeWithMessagesRange#365275f2 {X:Type} range:MessageRange query:!X = X;
// @method Invoke a method within a takeout session
// @param takeout_id Takeout session IDe
// @param query Query
invokeWithTakeout#aca9fd2e {X:Type} takeout_id:long query:!X = X;

// @method Send the verification code for login
// @param phone_number Phone number in international format
// @param api_id Application identifier (see App configuration)
// @param api_hash Application secret hash (see App configuration)
// @param settings Settings for the code type to send
auth.sendCode#a677244f phone_number:string api_id:int api_hash:string settings:CodeSettings = auth.SentCode;
// @method Registers a validated phone number in the system.
// @param phone_number Phone number in international format
// @param phone_code_hash SMS-message ID
// @param first_name New user first name
// @param last_name New user last name
auth.signUp#80eee427 phone_number:string phone_code_hash:string first_name:string last_name:string = auth.Authorization;
auth.signIn#bcd51581 phone_number:string phone_code_hash:string phone_code:string = auth.Authorization;
auth.logOut#5717da40 = Bool;
auth.resetAuthorizations#9fab0d1a = Bool;
auth.exportAuthorization#e5bfffcd dc_id:int = auth.ExportedAuthorization;
auth.importAuthorization#e3ef9613 id:int bytes:bytes = auth.Authorization;
auth.bindTempAuthKey#cdd42a05 perm_auth_key_id:long nonce:long expires_at:int encrypted_message:bytes = Bool;
auth.importBotAuthorization#67a3ff2c flags:int api_id:int api_hash:string bot_auth_token:string = auth.Authorization;
auth.checkPassword#d18b4d16 password:InputCheckPasswordSRP = auth.Authorization;
auth.requestPasswordRecovery#d897bc66 = auth.PasswordRecovery;
auth.recoverPassword#4ea56e92 code:string = auth.Authorization;
auth.resendCode#3ef1a9bf phone_number:string phone_code_hash:string = auth.SentCode;
auth.cancelCode#1f040578 phone_number:string phone_code_hash:string = Bool;
auth.dropTempAuthKeys#8e48a188 except_auth_keys:Vector<long> = Bool;
auth.exportLoginToken#b1b41517 api_id:int api_hash:string except_ids:Vector<int> = auth.LoginToken;
auth.importLoginToken#95ac5ce4 token:bytes = auth.LoginToken;
auth.acceptLoginToken#e894ad4d token:bytes = Authorization;

account.registerDevice#68976c6f flags:# no_muted:flags.0?true token_type:int token:string app_sandbox:Bool secret:bytes other_uids:Vector<int> = Bool;
account.unregisterDevice#3076c4bf token_type:int token:string other_uids:Vector<int> = Bool;
account.updateNotifySettings#84be5b93 peer:InputNotifyPeer settings:InputPeerNotifySettings = Bool;
account.getNotifySettings#12b3ad31 peer:InputNotifyPeer = PeerNotifySettings;
account.resetNotifySettings#db7e1747 = Bool;
account.updateProfile#78515775 flags:# first_name:flags.0?string last_name:flags.1?string about:flags.2?string = User;
account.updateStatus#6628562c offline:Bool = Bool;
account.getWallPapers#aabb1763 hash:int = account.WallPapers;
account.reportPeer#ae189d5f peer:InputPeer reason:ReportReason = Bool;
account.checkUsername#2714d86c username:string = Bool;
account.updateUsername#3e0bdd7c username:string = User;
account.getPrivacy#dadbc950 key:InputPrivacyKey = account.PrivacyRules;
account.setPrivacy#c9f81ce8 key:InputPrivacyKey rules:Vector<InputPrivacyRule> = account.PrivacyRules;
account.deleteAccount#418d4e0b reason:string = Bool;
account.getAccountTTL#8fc711d = AccountDaysTTL;
account.setAccountTTL#2442485e ttl:AccountDaysTTL = Bool;
account.sendChangePhoneCode#82574ae5 phone_number:string settings:CodeSettings = auth.SentCode;
account.changePhone#70c32edb phone_number:string phone_code_hash:string phone_code:string = User;
account.updateDeviceLocked#38df3532 period:int = Bool;
account.getAuthorizations#e320c158 = account.Authorizations;
account.resetAuthorization#df77f3bc hash:long = Bool;
account.getPassword#548a30f5 = account.Password;
account.getPasswordSettings#9cd4eaf9 password:InputCheckPasswordSRP = account.PasswordSettings;
account.updatePasswordSettings#a59b102f password:InputCheckPasswordSRP new_settings:account.PasswordInputSettings = Bool;
account.sendConfirmPhoneCode#1b3faa88 hash:string settings:CodeSettings = auth.SentCode;
account.confirmPhone#5f2178c3 phone_code_hash:string phone_code:string = Bool;
account.getTmpPassword#449e0b51 password:InputCheckPasswordSRP period:int = account.TmpPassword;
account.getWebAuthorizations#182e6d6f = account.WebAuthorizations;
account.resetWebAuthorization#2d01b9ef hash:long = Bool;
account.resetWebAuthorizations#682d2594 = Bool;
account.getAllSecureValues#b288bc7d = Vector<SecureValue>;
account.getSecureValue#73665bc2 types:Vector<SecureValueType> = Vector<SecureValue>;
account.saveSecureValue#899fe31d value:InputSecureValue secure_secret_id:long = SecureValue;
account.deleteSecureValue#b880bc4b types:Vector<SecureValueType> = Bool;
account.getAuthorizationForm#b86ba8e1 bot_id:int scope:string public_key:string = account.AuthorizationForm;
account.acceptAuthorization#e7027c94 bot_id:int scope:string public_key:string value_hashes:Vector<SecureValueHash> credentials:SecureCredentialsEncrypted = Bool;
account.sendVerifyPhoneCode#a5a356f9 phone_number:string settings:CodeSettings = auth.SentCode;
account.verifyPhone#4dd3a7f6 phone_number:string phone_code_hash:string phone_code:string = Bool;
account.sendVerifyEmailCode#7011509f email:string = account.SentEmailCode;
account.verifyEmail#ecba39db email:string code:string = Bool;
account.initTakeoutSession#f05b4804 flags:# contacts:flags.0?true message_users:flags.1?true message_chats:flags.2?true message_megagroups:flags.3?true message_channels:flags.4?true files:flags.5?true file_max_size:flags.5?int = account.Takeout;
account.finishTakeoutSession#1d2652ee flags:# success:flags.0?true = Bool;
account.confirmPasswordEmail#8fdf1920 code:string = Bool;
account.resendPasswordEmail#7a7f2a15 = Bool;
account.cancelPasswordEmail#c1cbd5b6 = Bool;
account.getContactSignUpNotification#9f07c728 = Bool;
account.setContactSignUpNotification#cff43f61 silent:Bool = Bool;
account.getNotifyExceptions#53577479 flags:# compare_sound:flags.1?true peer:flags.0?InputNotifyPeer = Updates;
account.getWallPaper#fc8ddbea wallpaper:InputWallPaper = WallPaper;
account.uploadWallPaper#dd853661 file:InputFile mime_type:string settings:WallPaperSettings = WallPaper;
account.saveWallPaper#6c5a5b37 wallpaper:InputWallPaper unsave:Bool settings:WallPaperSettings = Bool;
account.installWallPaper#feed5769 wallpaper:InputWallPaper settings:WallPaperSettings = Bool;
account.resetWallPapers#bb3b9804 = Bool;
account.getAutoDownloadSettings#56da0b3f = account.AutoDownloadSettings;
account.saveAutoDownloadSettings#76f36233 flags:# low:flags.0?true high:flags.1?true settings:AutoDownloadSettings = Bool;
account.uploadTheme#1c3db333 flags:# file:InputFile thumb:flags.0?InputFile file_name:string mime_type:string = Document;
account.createTheme#8432c21f flags:# slug:string title:string document:flags.2?InputDocument settings:flags.3?InputThemeSettings = Theme;
account.updateTheme#5cb367d5 flags:# format:string theme:InputTheme slug:flags.0?string title:flags.1?string document:flags.2?InputDocument settings:flags.3?InputThemeSettings = Theme;
account.saveTheme#f257106c theme:InputTheme unsave:Bool = Bool;
account.installTheme#7ae43737 flags:# dark:flags.0?true format:flags.1?string theme:flags.1?InputTheme = Bool;
account.getTheme#8d9d742b format:string theme:InputTheme document_id:long = Theme;
account.getThemes#285946f8 format:string hash:int = account.Themes;
account.setContentSettings#b574b16b flags:# sensitive_enabled:flags.0?true = Bool;
account.getContentSettings#8b9b4dae = account.ContentSettings;
account.getMultiWallPapers#65ad71dc wallpapers:Vector<InputWallPaper> = Vector<WallPaper>;
account.getGlobalPrivacySettings#eb2b4cf6 = GlobalPrivacySettings;
account.setGlobalPrivacySettings#1edaaac2 settings:GlobalPrivacySettings = GlobalPrivacySettings;

users.getUsers#d91a548 id:Vector<InputUser> = Vector<User>;
users.getFullUser#ca30a5b1 id:InputUser = UserFull;
users.setSecureValueErrors#90c894b5 id:InputUser errors:Vector<SecureValueError> = Bool;

// ***SPECIFIC RESPONSE***
contacts.getContactIDs#2caa4a42 hash:int = Vector<int>;
contacts.getStatuses#c4a353ee = Vector<ContactStatus>;
contacts.getContacts#c023849f hash:int = contacts.Contacts;
contacts.importContacts#2c800be5 contacts:Vector<InputContact> = contacts.ImportedContacts;
contacts.deleteContacts#96a0e00 id:Vector<InputUser> = Updates;
contacts.deleteByPhones#1013fd9e phones:Vector<string> = Bool;
contacts.block#68cc1411 id:InputPeer = Bool;
contacts.unblock#bea65d50 id:InputPeer = Bool;
contacts.getBlocked#f57c350f offset:int limit:int = contacts.Blocked;
contacts.search#11f812d8 q:string limit:int = contacts.Found;
contacts.resolveUsername#f93ccba3 username:string = contacts.ResolvedPeer;
contacts.getTopPeers#d4982db5 flags:# correspondents:flags.0?true bots_pm:flags.1?true bots_inline:flags.2?true phone_calls:flags.3?true forward_users:flags.4?true forward_chats:flags.5?true groups:flags.10?true channels:flags.15?true offset:int limit:int hash:int = contacts.TopPeers;
contacts.resetTopPeerRating#1ae373ac category:TopPeerCategory peer:InputPeer = Bool;
contacts.resetSaved#879537f1 = Bool;
contacts.getSaved#82f1e39f = Vector<SavedContact>;
contacts.toggleTopPeers#8514bdda enabled:Bool = Bool;
contacts.addContact#e8f463d0 flags:# add_phone_privacy_exception:flags.0?true id:InputUser first_name:string last_name:string phone:string = Updates;
contacts.acceptContact#f831a20f id:InputUser = Updates;
contacts.getLocated#d348bc44 flags:# background:flags.1?true geo_point:InputGeoPoint self_expires:flags.0?int = Updates;
contacts.blockFromReplies#29a8962c flags:# delete_message:flags.0?true delete_history:flags.1?true report_spam:flags.2?true msg_id:int = Updates;

messages.getMessages#63c66506 id:Vector<InputMessage> = messages.Messages;
messages.getDialogs#a0ee3b73 flags:# exclude_pinned:flags.0?true folder_id:flags.1?int offset_date:int offset_id:int offset_peer:InputPeer limit:int hash:int = messages.Dialogs;
messages.getHistory#dcbb8260 peer:InputPeer offset_id:int offset_date:int add_offset:int limit:int max_id:int min_id:int hash:int = messages.Messages;
messages.search#c352eec flags:# peer:InputPeer q:string from_id:flags.0?InputPeer top_msg_id:flags.1?int filter:MessagesFilter min_date:int max_date:int offset_id:int add_offset:int limit:int max_id:int min_id:int hash:int = messages.Messages;
messages.readHistory#e306d3a peer:InputPeer max_id:int = messages.AffectedMessages;
messages.deleteHistory#1c015b09 flags:# just_clear:flags.0?true revoke:flags.1?true peer:InputPeer max_id:int = messages.AffectedHistory;
messages.deleteMessages#e58e95d2 flags:# revoke:flags.0?true id:Vector<int> = messages.AffectedMessages;
messages.receivedMessages#5a954c0 max_id:int = Vector<ReceivedNotifyMessage>;
messages.setTyping#58943ee2 flags:# peer:InputPeer top_msg_id:flags.0?int action:SendMessageAction = Bool;
messages.sendMessage#520c3870 flags:# no_webpage:flags.1?true silent:flags.5?true background:flags.6?true clear_draft:flags.7?true peer:InputPeer reply_to_msg_id:flags.0?int message:string random_id:long reply_markup:flags.2?ReplyMarkup entities:flags.3?Vector<MessageEntity> schedule_date:flags.10?int = Updates;
messages.sendMedia#3491eba9 flags:# silent:flags.5?true background:flags.6?true clear_draft:flags.7?true peer:InputPeer reply_to_msg_id:flags.0?int media:InputMedia message:string random_id:long reply_markup:flags.2?ReplyMarkup entities:flags.3?Vector<MessageEntity> schedule_date:flags.10?int = Updates;
messages.forwardMessages#d9fee60e flags:# silent:flags.5?true background:flags.6?true with_my_score:flags.8?true from_peer:InputPeer id:Vector<int> random_id:Vector<long> to_peer:InputPeer schedule_date:flags.10?int = Updates;
messages.reportSpam#cf1592db peer:InputPeer = Bool;
messages.getPeerSettings#3672e09c peer:InputPeer = PeerSettings;
messages.report#bd82b658 peer:InputPeer id:Vector<int> reason:ReportReason = Bool;
messages.getChats#3c6aa187 id:Vector<int> = messages.Chats;
messages.getFullChat#3b831c66 chat_id:int = messages.ChatFull;
messages.editChatTitle#dc452855 chat_id:int title:string = Updates;
messages.editChatPhoto#ca4c79d8 chat_id:int photo:InputChatPhoto = Updates;
messages.addChatUser#f9a0aa09 chat_id:int user_id:InputUser fwd_limit:int = Updates;
messages.deleteChatUser#e0611f16 chat_id:int user_id:InputUser = Updates;
messages.createChat#9cb126e users:Vector<InputUser> title:string = Updates;
messages.getDhConfig#26cf8950 version:int random_length:int = messages.DhConfig;
messages.requestEncryption#f64daf43 user_id:InputUser random_id:int g_a:bytes = EncryptedChat;
messages.acceptEncryption#3dbc0415 peer:InputEncryptedChat g_b:bytes key_fingerprint:long = EncryptedChat;
messages.discardEncryption#edd923c5 chat_id:int = Bool;
messages.setEncryptedTyping#791451ed peer:InputEncryptedChat typing:Bool = Bool;
messages.readEncryptedHistory#7f4b690a peer:InputEncryptedChat max_date:int = Bool;
messages.sendEncrypted#44fa7a15 flags:# silent:flags.0?true peer:InputEncryptedChat random_id:long data:bytes = messages.SentEncryptedMessage;
messages.sendEncryptedFile#5559481d flags:# silent:flags.0?true peer:InputEncryptedChat random_id:long data:bytes file:InputEncryptedFile = messages.SentEncryptedMessage;
messages.sendEncryptedService#32d439a4 peer:InputEncryptedChat random_id:long data:bytes = messages.SentEncryptedMessage;
// ***SPECIFIC RESPONSE***
messages.receivedQueue#55a5bb66 max_qts:int = Vector<long>;
messages.reportEncryptedSpam#4b0c8c0f peer:InputEncryptedChat = Bool;
messages.readMessageContents#36a73f77 id:Vector<int> = messages.AffectedMessages;
messages.getStickers#43d4f2c emoticon:string hash:int = messages.Stickers;
messages.getAllStickers#1c9618b1 hash:int = messages.AllStickers;
messages.getWebPagePreview#8b68b0cc flags:# message:string entities:flags.3?Vector<MessageEntity> = MessageMedia;
messages.exportChatInvite#df7534c peer:InputPeer = ExportedChatInvite;
messages.checkChatInvite#3eadb1bb hash:string = ChatInvite;
messages.importChatInvite#6c50051c hash:string = Updates;
messages.getStickerSet#2619a90e stickerset:InputStickerSet = messages.StickerSet;
messages.installStickerSet#c78fe460 stickerset:InputStickerSet archived:Bool = messages.StickerSetInstallResult;
messages.uninstallStickerSet#f96e55de stickerset:InputStickerSet = Bool;
messages.startBot#e6df7378 bot:InputUser peer:InputPeer random_id:long start_param:string = Updates;
messages.getMessagesViews#5784d3e1 peer:InputPeer id:Vector<int> increment:Bool = messages.MessageViews;
messages.editChatAdmin#a9e69f2e chat_id:int user_id:InputUser is_admin:Bool = Bool;
messages.migrateChat#15a3b8e3 chat_id:int = Updates;
messages.searchGlobal#4bc6589a flags:# folder_id:flags.0?int q:string filter:MessagesFilter min_date:int max_date:int offset_rate:int offset_peer:InputPeer offset_id:int limit:int = messages.Messages;
messages.reorderStickerSets#78337739 flags:# masks:flags.0?true order:Vector<long> = Bool;
messages.getDocumentByHash#338e2464 sha256:bytes size:int mime_type:string = Document;
messages.getSavedGifs#83bf3d52 hash:int = messages.SavedGifs;
messages.saveGif#327a30cb id:InputDocument unsave:Bool = Bool;
messages.getInlineBotResults#514e999d flags:# bot:InputUser peer:InputPeer geo_point:flags.0?InputGeoPoint query:string offset:string = messages.BotResults;
messages.setInlineBotResults#eb5ea206 flags:# gallery:flags.0?true private:flags.1?true query_id:long results:Vector<InputBotInlineResult> cache_time:int next_offset:flags.2?string switch_pm:flags.3?InlineBotSwitchPM = Bool;
messages.sendInlineBotResult#220815b0 flags:# silent:flags.5?true background:flags.6?true clear_draft:flags.7?true hide_via:flags.11?true peer:InputPeer reply_to_msg_id:flags.0?int random_id:long query_id:long id:string schedule_date:flags.10?int = Updates;
messages.getMessageEditData#fda68d36 peer:InputPeer id:int = messages.MessageEditData;
messages.editMessage#48f71778 flags:# no_webpage:flags.1?true peer:InputPeer id:int message:flags.11?string media:flags.14?InputMedia reply_markup:flags.2?ReplyMarkup entities:flags.3?Vector<MessageEntity> schedule_date:flags.15?int = Updates;
messages.editInlineBotMessage#83557dba flags:# no_webpage:flags.1?true id:InputBotInlineMessageID message:flags.11?string media:flags.14?InputMedia reply_markup:flags.2?ReplyMarkup entities:flags.3?Vector<MessageEntity> = Bool;
messages.getBotCallbackAnswer#9342ca07 flags:# game:flags.1?true peer:InputPeer msg_id:int data:flags.0?bytes password:flags.2?InputCheckPasswordSRP = messages.BotCallbackAnswer;
messages.setBotCallbackAnswer#d58f130a flags:# alert:flags.1?true query_id:long message:flags.0?string url:flags.2?string cache_time:int = Bool;
messages.getPeerDialogs#e470bcfd peers:Vector<InputDialogPeer> = messages.PeerDialogs;
messages.saveDraft#bc39e14b flags:# no_webpage:flags.1?true reply_to_msg_id:flags.0?int peer:InputPeer message:string entities:flags.3?Vector<MessageEntity> = Bool;
messages.getAllDrafts#6a3f8d65 = Updates;
messages.getFeaturedStickers#2dacca4f hash:int = messages.FeaturedStickers;
messages.readFeaturedStickers#5b118126 id:Vector<long> = Bool;
messages.getRecentStickers#5ea192c9 flags:# attached:flags.0?true hash:int = messages.RecentStickers;
messages.saveRecentSticker#392718f8 flags:# attached:flags.0?true id:InputDocument unsave:Bool = Bool;
messages.clearRecentStickers#8999602d flags:# attached:flags.0?true = Bool;
messages.getArchivedStickers#57f17692 flags:# masks:flags.0?true offset_id:long limit:int = messages.ArchivedStickers;
messages.getMaskStickers#65b8c79f hash:int = messages.AllStickers;
messages.getAttachedStickers#cc5b67cc media:InputStickeredMedia = Vector<StickerSetCovered>;
messages.setGameScore#8ef8ecc0 flags:# edit_message:flags.0?true force:flags.1?true peer:InputPeer id:int user_id:InputUser score:int = Updates;
messages.setInlineGameScore#15ad9f64 flags:# edit_message:flags.0?true force:flags.1?true id:InputBotInlineMessageID user_id:InputUser score:int = Bool;
messages.getGameHighScores#e822649d peer:InputPeer id:int user_id:InputUser = messages.HighScores;
messages.getInlineGameHighScores#f635e1b id:InputBotInlineMessageID user_id:InputUser = messages.HighScores;
messages.getCommonChats#d0a48c4 user_id:InputUser max_id:int limit:int = messages.Chats;
messages.getAllChats#eba80ff0 except_ids:Vector<int> = messages.Chats;
messages.getWebPage#32ca8f91 url:string hash:int = WebPage;
messages.toggleDialogPin#a731e257 flags:# pinned:flags.0?true peer:InputDialogPeer = Bool;
messages.reorderPinnedDialogs#3b1adf37 flags:# force:flags.0?true folder_id:int order:Vector<InputDialogPeer> = Bool;
messages.getPinnedDialogs#d6b94df2 folder_id:int = messages.PeerDialogs;
messages.setBotShippingResults#e5f672fa flags:# query_id:long error:flags.0?string shipping_options:flags.1?Vector<ShippingOption> = Bool;
messages.setBotPrecheckoutResults#9c2dd95 flags:# success:flags.1?true query_id:long error:flags.0?string = Bool;
messages.uploadMedia#519bc2b1 peer:InputPeer media:InputMedia = MessageMedia;
messages.sendScreenshotNotification#c97df020 peer:InputPeer reply_to_msg_id:int random_id:long = Updates;
messages.getFavedStickers#21ce0b0e hash:int = messages.FavedStickers;
messages.faveSticker#b9ffc55b id:InputDocument unfave:Bool = Bool;
messages.getUnreadMentions#46578472 peer:InputPeer offset_id:int add_offset:int limit:int max_id:int min_id:int = messages.Messages;
messages.readMentions#f0189d3 peer:InputPeer = messages.AffectedHistory;
messages.getRecentLocations#bbc45b09 peer:InputPeer limit:int hash:int = messages.Messages;
messages.sendMultiMedia#cc0110cb flags:# silent:flags.5?true background:flags.6?true clear_draft:flags.7?true peer:InputPeer reply_to_msg_id:flags.0?int multi_media:Vector<InputSingleMedia> schedule_date:flags.10?int = Updates;
messages.uploadEncryptedFile#5057c497 peer:InputEncryptedChat file:InputEncryptedFile = EncryptedFile;
messages.searchStickerSets#c2b7d08b flags:# exclude_featured:flags.0?true q:string hash:int = messages.FoundStickerSets;
messages.getSplitRanges#1cff7e08 = Vector<MessageRange>;
messages.markDialogUnread#c286d98f flags:# unread:flags.0?true peer:InputDialogPeer = Bool;
messages.getDialogUnreadMarks#22e24e22 = Vector<DialogPeer>;
messages.clearAllDrafts#7e58ee9c = Bool;
messages.updatePinnedMessage#d2aaf7ec flags:# silent:flags.0?true unpin:flags.1?true pm_oneside:flags.2?true peer:InputPeer id:int = Updates;
messages.sendVote#10ea6184 peer:InputPeer msg_id:int options:Vector<bytes> = Updates;
messages.getPollResults#73bb643b peer:InputPeer msg_id:int = Updates;
messages.getOnlines#6e2be050 peer:InputPeer = ChatOnlines;
messages.getStatsURL#812c2ae6 flags:# dark:flags.0?true peer:InputPeer params:string = StatsURL;
messages.editChatAbout#def60797 peer:InputPeer about:string = Bool;
messages.editChatDefaultBannedRights#a5866b41 peer:InputPeer banned_rights:ChatBannedRights = Updates;
messages.getEmojiKeywords#35a0e062 lang_code:string = EmojiKeywordsDifference;
messages.getEmojiKeywordsDifference#1508b6af lang_code:string from_version:int = EmojiKeywordsDifference;
messages.getEmojiKeywordsLanguages#4e9963b2 lang_codes:Vector<string> = Vector<EmojiLanguage>;
messages.getEmojiURL#d5b10c26 lang_code:string = EmojiURL;
messages.getSearchCounters#732eef00 peer:InputPeer filters:Vector<MessagesFilter> = Vector<messages.SearchCounter>;
messages.requestUrlAuth#e33f5613 peer:InputPeer msg_id:int button_id:int = UrlAuthResult;
messages.acceptUrlAuth#f729ea98 flags:# write_allowed:flags.0?true peer:InputPeer msg_id:int button_id:int = UrlAuthResult;
messages.hidePeerSettingsBar#4facb138 peer:InputPeer = Bool;
messages.getScheduledHistory#e2c2685b peer:InputPeer hash:int = messages.Messages;
messages.getScheduledMessages#bdbb0464 peer:InputPeer id:Vector<int> = messages.Messages;
messages.sendScheduledMessages#bd38850a peer:InputPeer id:Vector<int> = Updates;
messages.deleteScheduledMessages#59ae2b16 peer:InputPeer id:Vector<int> = Updates;
messages.getPollVotes#b86e380e flags:# peer:InputPeer id:int option:flags.0?bytes offset:flags.1?string limit:int = messages.VotesList;
messages.toggleStickerSets#b5052fea flags:# uninstall:flags.0?true archive:flags.1?true unarchive:flags.2?true stickersets:Vector<InputStickerSet> = Bool;
messages.getDialogFilters#f19ed96d = Vector<DialogFilter>;
messages.getSuggestedDialogFilters#a29cd42c = Vector<DialogFilterSuggested>;
messages.updateDialogFilter#1ad4a04a flags:# id:int filter:flags.0?DialogFilter = Bool;
messages.updateDialogFiltersOrder#c563c1e4 order:Vector<int> = Bool;
messages.getOldFeaturedStickers#5fe7025b offset:int limit:int hash:int = messages.FeaturedStickers;
messages.getReplies#24b581ba peer:InputPeer msg_id:int offset_id:int offset_date:int add_offset:int limit:int max_id:int min_id:int hash:int = messages.Messages;
messages.getDiscussionMessage#446972fd peer:InputPeer msg_id:int = messages.DiscussionMessage;
messages.readDiscussion#f731a9f4 peer:InputPeer msg_id:int read_max_id:int = Bool;
messages.unpinAllMessages#f025bc8b peer:InputPeer = messages.AffectedHistory;

updates.getState#edd4882a = updates.State;
updates.getDifference#25939651 flags:# pts:int pts_total_limit:flags.0?int date:int qts:int = updates.Difference;
updates.getChannelDifference#3173d78 flags:# force:flags.0?true channel:InputChannel filter:ChannelMessagesFilter pts:int limit:int = updates.ChannelDifference;

photos.updateProfilePhoto#72d4742c id:InputPhoto = photos.Photo;
photos.uploadProfilePhoto#89f30f69 flags:# file:flags.0?InputFile video:flags.1?InputFile video_start_ts:flags.2?double = photos.Photo;
// ***SPECIFIC RESPONSE***
photos.deletePhotos#87cf7f2f id:Vector<InputPhoto> = Vector<long>;
photos.getUserPhotos#91cd32a8 user_id:InputUser offset:int max_id:long limit:int = photos.Photos;

upload.saveFilePart#b304a621 file_id:long file_part:int bytes:bytes = Bool;
upload.getFile#b15a9afc flags:# precise:flags.0?true cdn_supported:flags.1?true location:InputFileLocation offset:int limit:int = upload.File;
upload.saveBigFilePart#de7b673d file_id:long file_part:int file_total_parts:int bytes:bytes = Bool;
upload.getWebFile#24e6818d location:InputWebFileLocation offset:int limit:int = upload.WebFile;
upload.getCdnFile#2000bcc3 file_token:bytes offset:int limit:int = upload.CdnFile;
upload.reuploadCdnFile#9b2754a8 file_token:bytes request_token:bytes = Vector<FileHash>;
upload.getCdnFileHashes#4da54231 file_token:bytes offset:int = Vector<FileHash>;
upload.getFileHashes#c7025931 location:InputFileLocation offset:int = Vector<FileHash>;

help.getConfig#c4f9186b = Config;
help.getNearestDc#1fb33026 = NearestDc;
help.getAppUpdate#522d5a7d source:string = help.AppUpdate;
help.getInviteText#4d392343 = help.InviteText;
help.getSupport#9cdf08cd = help.Support;
help.getAppChangelog#9010ef6f prev_app_version:string = Updates;
help.setBotUpdatesStatus#ec22cfcd pending_updates_count:int message:string = Bool;
help.getCdnConfig#52029342 = CdnConfig;
help.getRecentMeUrls#3dc0f114 referer:string = help.RecentMeUrls;
help.getTermsOfServiceUpdate#2ca51fd1 = help.TermsOfServiceUpdate;
help.acceptTermsOfService#ee72f79a id:DataJSON = Bool;
help.getDeepLinkInfo#3fedc75f path:string = help.DeepLinkInfo;
help.getAppConfig#98914110 = JSONValue;
help.saveAppLog#6f02f748 events:Vector<InputAppEvent> = Bool;
help.getPassportConfig#c661ad08 hash:int = help.PassportConfig;
help.getSupportName#d360e72c = help.SupportName;
help.getUserInfo#38a08d3 user_id:InputUser = help.UserInfo;
help.editUserInfo#66b91b70 user_id:InputUser message:string entities:Vector<MessageEntity> = help.UserInfo;
help.getPromoData#c0977421 = help.PromoData;
help.hidePromoData#1e251c95 peer:InputPeer = Bool;
help.dismissSuggestion#77fa99f suggestion:string = Bool;
help.getCountriesList#735787a8 lang_code:string hash:int = help.CountriesList;

channels.readHistory#cc104937 channel:InputChannel max_id:int = Bool;
channels.deleteMessages#84c1fd4e channel:InputChannel id:Vector<int> = messages.AffectedMessages;
channels.deleteUserHistory#d10dd71b channel:InputChannel user_id:InputUser = messages.AffectedHistory;
channels.reportSpam#fe087810 channel:InputChannel user_id:InputUser id:Vector<int> = Bool;
channels.getMessages#ad8c9a23 channel:InputChannel id:Vector<InputMessage> = messages.Messages;
// @method Get the participants of a channel
// @param channel Channel
// @param filter Which participant types to fetch
channels.getParticipants#123e05e9 channel:InputChannel filter:ChannelParticipantsFilter offset:int limit:int hash:int = channels.ChannelParticipants;
channels.getParticipant#546dd7a6 channel:InputChannel user_id:InputUser = channels.ChannelParticipant;
channels.getChannels#a7f6bbb id:Vector<InputChannel> = messages.Chats;
channels.getFullChannel#8736a09 channel:InputChannel = messages.ChatFull;
channels.createChannel#3d5fb10f flags:# broadcast:flags.0?true megagroup:flags.1?true title:string about:string geo_point:flags.2?InputGeoPoint address:flags.2?string = Updates;
channels.editAdmin#d33c8902 channel:InputChannel user_id:InputUser admin_rights:ChatAdminRights rank:string = Updates;
channels.editTitle#566decd0 channel:InputChannel title:string = Updates;
channels.editPhoto#f12e57c9 channel:InputChannel photo:InputChatPhoto = Updates;
channels.checkUsername#10e6bd2c channel:InputChannel username:string = Bool;
channels.updateUsername#3514b3de channel:InputChannel username:string = Bool;
channels.joinChannel#24b524c5 channel:InputChannel = Updates;
channels.leaveChannel#f836aa95 channel:InputChannel = Updates;
channels.inviteToChannel#199f3a6c channel:InputChannel users:Vector<InputUser> = Updates;
channels.deleteChannel#c0111fe3 channel:InputChannel = Updates;
channels.exportMessageLink#e63fadeb flags:# grouped:flags.0?true thread:flags.1?true channel:InputChannel id:int = ExportedMessageLink;
channels.toggleSignatures#1f69b606 channel:InputChannel enabled:Bool = Updates;
channels.getAdminedPublicChannels#f8b036af flags:# by_location:flags.0?true check_limit:flags.1?true = messages.Chats;
channels.editBanned#72796912 channel:InputChannel user_id:InputUser banned_rights:ChatBannedRights = Updates;
channels.getAdminLog#33ddf480 flags:# channel:InputChannel q:string events_filter:flags.0?ChannelAdminLogEventsFilter admins:flags.1?Vector<InputUser> max_id:long min_id:long limit:int = channels.AdminLogResults;
channels.setStickers#ea8ca4f9 channel:InputChannel stickerset:InputStickerSet = Bool;
channels.readMessageContents#eab5dc38 channel:InputChannel id:Vector<int> = Bool;
channels.deleteHistory#af369d42 channel:InputChannel max_id:int = Bool;
channels.togglePreHistoryHidden#eabbb94c channel:InputChannel enabled:Bool = Updates;
channels.getLeftChannels#8341ecc0 offset:int = messages.Chats;
channels.getGroupsForDiscussion#f5dad378 = messages.Chats;
channels.setDiscussionGroup#40582bb2 broadcast:InputChannel group:InputChannel = Bool;
channels.editCreator#8f38cd1f channel:InputChannel user_id:InputUser password:InputCheckPasswordSRP = Updates;
channels.editLocation#58e63f6d channel:InputChannel geo_point:InputGeoPoint address:string = Bool;
channels.toggleSlowMode#edd49ef0 channel:InputChannel seconds:int = Updates;
channels.getInactiveChannels#11e831ee = messages.InactiveChats;

bots.sendCustomRequest#aa2769ed custom_method:string params:DataJSON = DataJSON;
bots.answerWebhookJSONQuery#e6213f4d query_id:long data:DataJSON = Bool;
bots.setBotCommands#805d46f6 commands:Vector<BotCommand> = Bool;

payments.getPaymentForm#99f09745 msg_id:int = payments.PaymentForm;
payments.getPaymentReceipt#a092a980 msg_id:int = payments.PaymentReceipt;
payments.validateRequestedInfo#770a8e74 flags:# save:flags.0?true msg_id:int info:PaymentRequestedInfo = payments.ValidatedRequestedInfo;
payments.sendPaymentForm#2b8879b3 flags:# msg_id:int requested_info_id:flags.0?string shipping_option_id:flags.1?string credentials:InputPaymentCredentials = payments.PaymentResult;
payments.getSavedInfo#227d824b = payments.SavedInfo;
payments.clearSavedInfo#d83d70c1 flags:# credentials:flags.0?true info:flags.1?true = Bool;
payments.getBankCardData#2e79d779 number:string = payments.BankCardData;

stickers.createStickerSet#f1036780 flags:# masks:flags.0?true animated:flags.1?true user_id:InputUser title:string short_name:string thumb:flags.2?InputDocument stickers:Vector<InputStickerSetItem> = messages.StickerSet;
stickers.removeStickerFromSet#f7760f51 sticker:InputDocument = messages.StickerSet;
stickers.changeStickerPosition#ffb6d4ca sticker:InputDocument position:int = messages.StickerSet;
stickers.addStickerToSet#8653febe stickerset:InputStickerSet sticker:InputStickerSetItem = messages.StickerSet;
stickers.setStickerSetThumb#9a364e30 stickerset:InputStickerSet thumb:InputDocument = messages.StickerSet;

phone.getCallConfig#55451fa9 = DataJSON;
phone.requestCall#42ff96ed flags:# video:flags.0?true user_id:InputUser random_id:int g_a_hash:bytes protocol:PhoneCallProtocol = phone.PhoneCall;
phone.acceptCall#3bd2b4a0 peer:InputPhoneCall g_b:bytes protocol:PhoneCallProtocol = phone.PhoneCall;
phone.confirmCall#2efe1722 peer:InputPhoneCall g_a:bytes key_fingerprint:long protocol:PhoneCallProtocol = phone.PhoneCall;
phone.receivedCall#17d54f61 peer:InputPhoneCall = Bool;
phone.discardCall#b2cbc1c0 flags:# video:flags.0?true peer:InputPhoneCall duration:int reason:PhoneCallDiscardReason connection_id:long = Updates;
phone.setCallRating#59ead627 flags:# user_initiative:flags.0?true peer:InputPhoneCall rating:int comment:string = Updates;
phone.saveCallDebug#277add7e peer:InputPhoneCall debug:DataJSON = Bool;
phone.sendSignalingData#ff7a9383 peer:InputPhoneCall data:bytes = Bool;

langpack.getLangPack#f2f2330a lang_pack:string lang_code:string = LangPackDifference;
langpack.getStrings#efea3803 lang_pack:string lang_code:string keys:Vector<string> = Vector<LangPackString>;
langpack.getDifference#cd984aa5 lang_pack:string lang_code:string from_version:int = LangPackDifference;
langpack.getLanguages#42c6978f lang_pack:string = Vector<LangPackLanguage>;
langpack.getLanguage#6a596502 lang_pack:string lang_code:string = LangPackLanguage;

folders.editPeerFolders#6847d0ab folder_peers:Vector<InputFolderPeer> = Updates;
folders.deleteFolder#1c295881 folder_id:int = Updates;

stats.getBroadcastStats#ab42441a flags:# dark:flags.0?true channel:InputChannel = stats.BroadcastStats;
stats.loadAsyncGraph#621d5fa0 flags:# token:string x:flags.0?long = StatsGraph;
stats.getMegagroupStats#dcdf8607 flags:# dark:flags.0?true channel:InputChannel = stats.MegagroupStats;
stats.getMessagePublicForwards#5630281b channel:InputChannel msg_id:int offset_rate:int offset_peer:InputPeer offset_id:int limit:int = messages.Messages;
stats.getMessageStats#b6e0a3f5 flags:# dark:flags.0?true channel:InputChannel msg_id:int = stats.MessageStats;