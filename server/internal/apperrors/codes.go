package apperrors

const (
	CodeInternal              = "INTERNAL_ERROR"
	CodeDependencyUnavailable = "DEPENDENCY_UNAVAILABLE"
)

// Plugin auth
const (
	CodePluginAuthV1HandlerKeyMissing = "PLUGIN_AUTH_V1_HANDLER_KEY_MISSING"
	MsgPluginAuthV1HandlerKeyMissing  = "Plugin API key is required."

	CodePluginAuthV1HandlerKeyInvalid = "PLUGIN_AUTH_V1_HANDLER_KEY_INVALID"
	MsgPluginAuthV1HandlerKeyInvalid  = "Plugin API key is invalid."
)

// GET /ready
const (
	CodeReadyGetHandlerSqlGetterFailed = "READY_GET_HANDLER_SQL_GETTER_FAILED"
	MsgReadyGetHandlerSqlGetterFailed  = "Database connection handle is not available."

	CodeReadyGetHandlerDatabasePingFailed = "READY_GET_HANDLER_DATABASE_PING_FAILED"
	MsgReadyGetHandlerDatabasePingFailed  = "Database is not reachable."
)

// GET /v1/internal/whitelist/{minecraftUuid}
const (
	CodeWhitelistGetV1HandlerUUIDInvalid = "WHITELIST_GET_V1_HANDLER_UUID_INVALID"
	MsgWhitelistGetV1HandlerUUIDInvalid  = "minecraftUuid must be a valid UUID."

	CodeWhitelistGetV1ServiceCheckFailed = "WHITELIST_GET_V1_SERVICE_CHECK_FAILED"
	MsgWhitelistGetV1ServiceCheckFailed  = "Could not check whitelist."
)

// POST /v1/internal/players/upsert
const (
	CodePlayerUpsertV1HandlerBodyInvalid = "PLAYER_UPSERT_V1_HANDLER_BODY_INVALID"
	MsgPlayerUpsertV1HandlerBodyInvalid  = "Request body must include minecraftUuid, username, and serverSlug."

	CodePlayerUpsertV1ServiceFailed = "PLAYER_UPSERT_V1_SERVICE_FAILED"
	MsgPlayerUpsertV1ServiceFailed  = "Could not upsert player."
)

// GET /v1/internal/whitelist/bedrock/{xuid}
const (
	CodeBedrockWhitelistGetV1HandlerXUIDInvalid = "BEDROCK_WHITELIST_GET_V1_HANDLER_XUID_INVALID"
	MsgBedrockWhitelistGetV1HandlerXUIDInvalid  = "xuid must be a non-empty numeric Xbox user id."

	CodeBedrockWhitelistGetV1ServiceCheckFailed = "BEDROCK_WHITELIST_GET_V1_SERVICE_CHECK_FAILED"
	MsgBedrockWhitelistGetV1ServiceCheckFailed  = "Could not check bedrock whitelist."
)

// POST /v1/internal/players/bedrock/upsert
const (
	CodeBedrockPlayerUpsertV1HandlerBodyInvalid = "BEDROCK_PLAYER_UPSERT_V1_HANDLER_BODY_INVALID"
	MsgBedrockPlayerUpsertV1HandlerBodyInvalid  = "Request body must include xuid, username, and serverSlug."

	CodeBedrockPlayerUpsertV1ServiceFailed = "BEDROCK_PLAYER_UPSERT_V1_SERVICE_FAILED"
	MsgBedrockPlayerUpsertV1ServiceFailed  = "Could not upsert bedrock player."
)

// POST /v1/internal/presence/*
const (
	CodePresenceOnlineV1HandlerBodyInvalid = "PRESENCE_ONLINE_V1_HANDLER_BODY_INVALID"
	MsgPresenceOnlineV1HandlerBodyInvalid  = "Request body must include playerId, username, and serverSlug."

	CodePresenceOnlineV1ServiceFailed = "PRESENCE_ONLINE_V1_SERVICE_FAILED"
	MsgPresenceOnlineV1ServiceFailed  = "Could not mark player online."

	CodePresenceOfflineV1HandlerBodyInvalid = "PRESENCE_OFFLINE_V1_HANDLER_BODY_INVALID"
	MsgPresenceOfflineV1HandlerBodyInvalid  = "Request body must include playerId and serverSlug."

	CodePresenceOfflineV1ServiceFailed = "PRESENCE_OFFLINE_V1_SERVICE_FAILED"
	MsgPresenceOfflineV1ServiceFailed  = "Could not mark player offline."

	CodePresenceHeartbeatV1HandlerBodyInvalid = "PRESENCE_HEARTBEAT_V1_HANDLER_BODY_INVALID"
	MsgPresenceHeartbeatV1HandlerBodyInvalid  = "Request body must include playerId and serverSlug."

	CodePresenceHeartbeatV1ServiceFailed = "PRESENCE_HEARTBEAT_V1_SERVICE_FAILED"
	MsgPresenceHeartbeatV1ServiceFailed  = "Could not refresh presence heartbeat."
)

// GET /v1/presence/*
const (
	CodePresenceOverviewV1ServiceFailed = "PRESENCE_OVERVIEW_V1_SERVICE_FAILED"
	MsgPresenceOverviewV1ServiceFailed  = "Could not load presence overview."

	CodePresenceServerV1HandlerSlugInvalid = "PRESENCE_SERVER_V1_HANDLER_SLUG_INVALID"
	MsgPresenceServerV1HandlerSlugInvalid  = "server slug is required."

	CodePresenceServerV1ServiceFailed = "PRESENCE_SERVER_V1_SERVICE_FAILED"
	MsgPresenceServerV1ServiceFailed  = "Could not load server presence."

	CodePresenceGuildV1ServiceFailed = "PRESENCE_GUILD_V1_SERVICE_FAILED"
	MsgPresenceGuildV1ServiceFailed  = "Could not load guild presence."
)

// POST /v1/internal/stats/ingest
const (
	CodeStatsIngestV1HandlerBodyInvalid = "STATS_INGEST_V1_HANDLER_BODY_INVALID"
	MsgStatsIngestV1HandlerBodyInvalid  = "Request body must include playerId, serverSlug, and sessionSeconds or mobKills."

	CodeStatsIngestV1ServiceFailed = "STATS_INGEST_V1_SERVICE_FAILED"
	MsgStatsIngestV1ServiceFailed  = "Could not ingest player stats."
)

// GET /v1/players/{id}/stats
const (
	CodeStatsGetV1ServiceFailed = "STATS_GET_V1_SERVICE_FAILED"
	MsgStatsGetV1ServiceFailed  = "Could not load player stats."
)

// GET /v1/internal/players/{id}/hud
const (
	CodeStatsHUDV1ServiceFailed = "STATS_HUD_V1_SERVICE_FAILED"
	MsgStatsHUDV1ServiceFailed  = "Could not load player HUD."
)

// POST /v1/internal/invites
const (
	CodeInvitePostInternalV1HandlerBodyInvalid = "INVITE_POST_INTERNAL_V1_HANDLER_BODY_INVALID"
	MsgInvitePostInternalV1HandlerBodyInvalid  = "Request body must include sponsorUuid and targetUsername."

	CodeInvitePostInternalV1ServiceSponsorNotFound = "INVITE_POST_INTERNAL_V1_SERVICE_SPONSOR_NOT_FOUND"
	MsgInvitePostInternalV1ServiceSponsorNotFound  = "Sponsor player not found."

	CodeInvitePostInternalV1ServiceSponsorBanned = "INVITE_POST_INTERNAL_V1_SERVICE_SPONSOR_BANNED"
	MsgInvitePostInternalV1ServiceSponsorBanned  = "Banned players cannot send invites."

	CodeInvitePostInternalV1ServiceSponsorProbation = "INVITE_POST_INTERNAL_V1_SERVICE_SPONSOR_PROBATION"
	MsgInvitePostInternalV1ServiceSponsorProbation  = "Players on probation cannot send invites."

	CodeInvitePostInternalV1ServicePendingExists = "INVITE_POST_INTERNAL_V1_SERVICE_PENDING_EXISTS"
	MsgInvitePostInternalV1ServicePendingExists  = "A pending invite already exists for this username."

	CodeInvitePostInternalV1ServiceCreateFailed = "INVITE_POST_INTERNAL_V1_SERVICE_CREATE_FAILED"
	MsgInvitePostInternalV1ServiceCreateFailed  = "Could not create invite."
)

// GET /v1/players/{id}
const (
	CodePlayerGetV1HandlerIDInvalid = "PLAYER_GET_V1_HANDLER_ID_INVALID"
	MsgPlayerGetV1HandlerIDInvalid  = "Player id must be a valid UUID."

	CodePlayerGetV1ServiceNotFound = "PLAYER_GET_V1_SERVICE_NOT_FOUND"
	MsgPlayerGetV1ServiceNotFound  = "Player not found."

	CodePlayerGetV1ServiceLoadFailed = "PLAYER_GET_V1_SERVICE_LOAD_FAILED"
	MsgPlayerGetV1ServiceLoadFailed  = "Could not load player."
)

// GET /v1/lookup/players/minecraft/{minecraftUuid}
const (
	CodePlayerGetMinecraftV1HandlerUUIDInvalid = "PLAYER_GET_MINECRAFT_V1_HANDLER_UUID_INVALID"
	MsgPlayerGetMinecraftV1HandlerUUIDInvalid  = "minecraftUuid must be a valid UUID."

	CodePlayerGetMinecraftV1ServiceNotFound = "PLAYER_GET_MINECRAFT_V1_SERVICE_NOT_FOUND"
	MsgPlayerGetMinecraftV1ServiceNotFound  = "Player not found."

	CodePlayerGetMinecraftV1ServiceLoadFailed = "PLAYER_GET_MINECRAFT_V1_SERVICE_LOAD_FAILED"
	MsgPlayerGetMinecraftV1ServiceLoadFailed  = "Could not load player."
)

// GET /v1/players/{id}/invites
const (
	CodePlayerInvitesListV1HandlerIDInvalid = "PLAYER_INVITES_LIST_V1_HANDLER_ID_INVALID"
	MsgPlayerInvitesListV1HandlerIDInvalid  = "Player id must be a valid UUID."

	CodePlayerInvitesListV1ServiceNotFound = "PLAYER_INVITES_LIST_V1_SERVICE_NOT_FOUND"
	MsgPlayerInvitesListV1ServiceNotFound  = "Player not found."

	CodePlayerInvitesListV1ServiceListFailed = "PLAYER_INVITES_LIST_V1_SERVICE_LIST_FAILED"
	MsgPlayerInvitesListV1ServiceListFailed  = "Could not list invites."
)

// GET /v1/invites/{code}
const (
	CodeInviteGetV1HandlerCodeEmpty = "INVITE_GET_V1_HANDLER_CODE_EMPTY"
	MsgInviteGetV1HandlerCodeEmpty  = "Invite code is required."

	CodeInviteGetV1ServiceNotFound = "INVITE_GET_V1_SERVICE_NOT_FOUND"
	MsgInviteGetV1ServiceNotFound  = "Invite not found."

	CodeInviteGetV1ServiceLoadFailed = "INVITE_GET_V1_SERVICE_LOAD_FAILED"
	MsgInviteGetV1ServiceLoadFailed  = "Could not load invite."
)

// Guilds
const (
	CodeGuildPostV1HandlerBodyInvalid = "GUILD_POST_V1_HANDLER_BODY_INVALID"
	MsgGuildPostV1HandlerBodyInvalid  = "Request body must include leaderUuid and name."

	CodeGuildPostV1ServiceNameInvalid = "GUILD_POST_V1_SERVICE_NAME_INVALID"
	MsgGuildPostV1ServiceNameInvalid  = "Guild name must be between 1 and 64 characters."

	CodeGuildPostV1ServiceLeaderNotFound = "GUILD_POST_V1_SERVICE_LEADER_NOT_FOUND"
	MsgGuildPostV1ServiceLeaderNotFound  = "Leader player not found."

	CodeGuildPostV1ServiceLeaderBanned = "GUILD_POST_V1_SERVICE_LEADER_BANNED"
	MsgGuildPostV1ServiceLeaderBanned  = "Banned players cannot create guilds."

	CodeGuildPostV1ServiceLeaderProbation = "GUILD_POST_V1_SERVICE_LEADER_PROBATION"
	MsgGuildPostV1ServiceLeaderProbation  = "Players on probation cannot create guilds."

	CodeGuildPostV1ServiceAlreadyMember = "GUILD_POST_V1_SERVICE_ALREADY_MEMBER"
	MsgGuildPostV1ServiceAlreadyMember  = "Player already belongs to a guild."

	CodeGuildPostV1ServiceCreateFailed = "GUILD_POST_V1_SERVICE_CREATE_FAILED"
	MsgGuildPostV1ServiceCreateFailed  = "Could not create guild."

	CodeGuildGetV1HandlerIDInvalid = "GUILD_GET_V1_HANDLER_ID_INVALID"
	MsgGuildGetV1HandlerIDInvalid  = "Guild id must be a valid UUID."

	CodeGuildGetV1ServiceNotFound = "GUILD_GET_V1_SERVICE_NOT_FOUND"
	MsgGuildGetV1ServiceNotFound  = "Guild not found."

	CodeGuildGetV1ServiceLoadFailed = "GUILD_GET_V1_SERVICE_LOAD_FAILED"
	MsgGuildGetV1ServiceLoadFailed  = "Could not load guild."

	CodeGuildListV1ServiceListFailed = "GUILD_LIST_V1_SERVICE_LIST_FAILED"
	MsgGuildListV1ServiceListFailed  = "Could not list guilds."

	CodeGuildMembersV1ServiceGuildNotFound = "GUILD_MEMBERS_V1_SERVICE_GUILD_NOT_FOUND"
	MsgGuildMembersV1ServiceGuildNotFound  = "Guild not found."

	CodeGuildMembersV1ServiceListFailed = "GUILD_MEMBERS_V1_SERVICE_LIST_FAILED"
	MsgGuildMembersV1ServiceListFailed  = "Could not list guild members."

	CodeGuildJoinV1ServicePlayerNotFound = "GUILD_JOIN_V1_SERVICE_PLAYER_NOT_FOUND"
	MsgGuildJoinV1ServicePlayerNotFound  = "Player not found."

	CodeGuildJoinV1ServiceGuildNotFound = "GUILD_JOIN_V1_SERVICE_GUILD_NOT_FOUND"
	MsgGuildJoinV1ServiceGuildNotFound  = "Guild not found."

	CodeGuildJoinV1ServiceProbation = "GUILD_JOIN_V1_SERVICE_PROBATION"
	MsgGuildJoinV1ServiceProbation  = "Players on probation cannot join guilds."

	CodeGuildJoinV1ServiceAlreadyMember = "GUILD_JOIN_V1_SERVICE_ALREADY_MEMBER"
	MsgGuildJoinV1ServiceAlreadyMember  = "Player already belongs to a guild."

	CodeGuildJoinV1ServiceFailed = "GUILD_JOIN_V1_SERVICE_FAILED"
	MsgGuildJoinV1ServiceFailed  = "Could not join guild."

	CodeGuildLeaveV1ServicePlayerNotFound = "GUILD_LEAVE_V1_SERVICE_PLAYER_NOT_FOUND"
	MsgGuildLeaveV1ServicePlayerNotFound  = "Player not found."

	CodeGuildLeaveV1ServiceGuildNotFound = "GUILD_LEAVE_V1_SERVICE_GUILD_NOT_FOUND"
	MsgGuildLeaveV1ServiceGuildNotFound  = "Guild not found."

	CodeGuildLeaveV1ServiceLeaderCannotLeave = "GUILD_LEAVE_V1_SERVICE_LEADER_CANNOT_LEAVE"
	MsgGuildLeaveV1ServiceLeaderCannotLeave  = "Guild leader cannot leave; transfer leadership first."

	CodeGuildLeaveV1ServiceNotMember = "GUILD_LEAVE_V1_SERVICE_NOT_MEMBER"
	MsgGuildLeaveV1ServiceNotMember  = "Player is not a member of this guild."

	CodeGuildLeaveV1ServiceFailed = "GUILD_LEAVE_V1_SERVICE_FAILED"
	MsgGuildLeaveV1ServiceFailed  = "Could not leave guild."
)

// Trust
const (
	CodeTrustEventV1HandlerBodyInvalid = "TRUST_EVENT_V1_HANDLER_BODY_INVALID"
	MsgTrustEventV1HandlerBodyInvalid  = "Request body must include playerId and eventType."

	CodeTrustEventV1ServicePlayerNotFound = "TRUST_EVENT_V1_SERVICE_PLAYER_NOT_FOUND"
	MsgTrustEventV1ServicePlayerNotFound  = "Player not found."

	CodeTrustEventV1ServiceRecordFailed = "TRUST_EVENT_V1_SERVICE_RECORD_FAILED"
	MsgTrustEventV1ServiceRecordFailed  = "Could not record trust event."

	CodeTrustListV1ServicePlayerNotFound = "TRUST_LIST_V1_SERVICE_PLAYER_NOT_FOUND"
	MsgTrustListV1ServicePlayerNotFound  = "Player not found."

	CodeTrustListV1ServiceListFailed = "TRUST_LIST_V1_SERVICE_LIST_FAILED"
	MsgTrustListV1ServiceListFailed  = "Could not list trust events."

	CodeTrustTreeV1ServicePlayerNotFound = "TRUST_TREE_V1_SERVICE_PLAYER_NOT_FOUND"
	MsgTrustTreeV1ServicePlayerNotFound  = "Player not found."

	CodeTrustTreeV1ServiceLoadFailed = "TRUST_TREE_V1_SERVICE_LOAD_FAILED"
	MsgTrustTreeV1ServiceLoadFailed  = "Could not load sponsor tree."
)

// Cities
const (
	CodeCityPostV1HandlerBodyInvalid = "CITY_POST_V1_HANDLER_BODY_INVALID"
	MsgCityPostV1HandlerBodyInvalid  = "Request body must include founderUuid, name, and serverSlug."

	CodeCityPostV1ServiceNameInvalid = "CITY_POST_V1_SERVICE_NAME_INVALID"
	MsgCityPostV1ServiceNameInvalid  = "City name must be between 1 and 64 characters."

	CodeCityPostV1ServiceFounderNotFound = "CITY_POST_V1_SERVICE_FOUNDER_NOT_FOUND"
	MsgCityPostV1ServiceFounderNotFound  = "Founder player not found."

	CodeCityPostV1ServiceFounderRestricted = "CITY_POST_V1_SERVICE_FOUNDER_RESTRICTED"
	MsgCityPostV1ServiceFounderRestricted  = "Founder cannot create cities in current status."

	CodeCityPostV1ServiceGuildNotFound = "CITY_POST_V1_SERVICE_GUILD_NOT_FOUND"
	MsgCityPostV1ServiceGuildNotFound  = "Guild not found."

	CodeCityPostV1ServiceCreateFailed = "CITY_POST_V1_SERVICE_CREATE_FAILED"
	MsgCityPostV1ServiceCreateFailed  = "Could not create city."

	CodeCityGetV1HandlerIDInvalid = "CITY_GET_V1_HANDLER_ID_INVALID"
	MsgCityGetV1HandlerIDInvalid  = "City id must be a valid UUID."

	CodeCityGetV1ServiceNotFound = "CITY_GET_V1_SERVICE_NOT_FOUND"
	MsgCityGetV1ServiceNotFound  = "City not found."

	CodeCityGetV1ServiceLoadFailed = "CITY_GET_V1_SERVICE_LOAD_FAILED"
	MsgCityGetV1ServiceLoadFailed  = "Could not load city."

	CodeCityListV1ServiceListFailed = "CITY_LIST_V1_SERVICE_LIST_FAILED"
	MsgCityListV1ServiceListFailed  = "Could not list cities."
)

// Claims
const (
	CodeClaimPostV1HandlerBodyInvalid = "CLAIM_POST_V1_HANDLER_BODY_INVALID"
	MsgClaimPostV1HandlerBodyInvalid  = "Request body must include ownerUuid, serverSlug, and bounds."

	CodeClaimPostV1ServiceOwnerNotFound = "CLAIM_POST_V1_SERVICE_OWNER_NOT_FOUND"
	MsgClaimPostV1ServiceOwnerNotFound  = "Owner player not found."

	CodeClaimPostV1ServiceProbation = "CLAIM_POST_V1_SERVICE_PROBATION"
	MsgClaimPostV1ServiceProbation  = "Players on probation cannot claim land."

	CodeClaimPostV1ServiceBanned = "CLAIM_POST_V1_SERVICE_BANNED"
	MsgClaimPostV1ServiceBanned  = "Banned players cannot claim land."

	CodeClaimPostV1ServiceBoundsInvalid = "CLAIM_POST_V1_SERVICE_BOUNDS_INVALID"
	MsgClaimPostV1ServiceBoundsInvalid  = "Claim bounds are invalid."

	CodeClaimPostV1ServiceAreaLimit = "CLAIM_POST_V1_SERVICE_AREA_LIMIT"
	MsgClaimPostV1ServiceAreaLimit  = "Claim would exceed territorial area limit."

	CodeClaimPostV1ServiceZoneInvalid = "CLAIM_POST_V1_SERVICE_ZONE_INVALID"
	MsgClaimPostV1ServiceZoneInvalid  = "Zone type must be urban, rural, industrial, or historic."

	CodeClaimPostV1ServiceCityNotFound = "CLAIM_POST_V1_SERVICE_CITY_NOT_FOUND"
	MsgClaimPostV1ServiceCityNotFound  = "City not found."

	CodeClaimPostV1ServiceOverlap = "CLAIM_POST_V1_SERVICE_OVERLAP"
	MsgClaimPostV1ServiceOverlap  = "Claim overlaps an existing claim."

	CodeClaimPostV1ServiceCreateFailed = "CLAIM_POST_V1_SERVICE_CREATE_FAILED"
	MsgClaimPostV1ServiceCreateFailed  = "Could not create claim."

	CodeClaimDeleteV1HandlerIDInvalid = "CLAIM_DELETE_V1_HANDLER_ID_INVALID"
	MsgClaimDeleteV1HandlerIDInvalid  = "Claim id must be a valid UUID."

	CodeClaimDeleteV1ServiceNotFound = "CLAIM_DELETE_V1_SERVICE_NOT_FOUND"
	MsgClaimDeleteV1ServiceNotFound  = "Claim not found."

	CodeClaimDeleteV1ServiceOwnerNotFound = "CLAIM_DELETE_V1_SERVICE_OWNER_NOT_FOUND"
	MsgClaimDeleteV1ServiceOwnerNotFound  = "Owner player not found."

	CodeClaimDeleteV1ServiceNotOwner = "CLAIM_DELETE_V1_SERVICE_NOT_OWNER"
	MsgClaimDeleteV1ServiceNotOwner  = "Only the claim owner can delete it."

	CodeClaimDeleteV1ServiceFailed = "CLAIM_DELETE_V1_SERVICE_FAILED"
	MsgClaimDeleteV1ServiceFailed  = "Could not delete claim."

	CodeClaimGetV1HandlerIDInvalid = "CLAIM_GET_V1_HANDLER_ID_INVALID"
	MsgClaimGetV1HandlerIDInvalid  = "Claim id must be a valid UUID."

	CodeClaimGetV1ServiceNotFound = "CLAIM_GET_V1_SERVICE_NOT_FOUND"
	MsgClaimGetV1ServiceNotFound  = "Claim not found."

	CodeClaimGetV1ServiceLoadFailed = "CLAIM_GET_V1_SERVICE_LOAD_FAILED"
	MsgClaimGetV1ServiceLoadFailed  = "Could not load claim."

	CodeClaimListCityV1HandlerIDInvalid = "CLAIM_LIST_CITY_V1_HANDLER_ID_INVALID"
	MsgClaimListCityV1HandlerIDInvalid  = "City id must be a valid UUID."

	CodeClaimListCityV1ServiceCityNotFound = "CLAIM_LIST_CITY_V1_SERVICE_CITY_NOT_FOUND"
	MsgClaimListCityV1ServiceCityNotFound  = "City not found."

	CodeClaimListCityV1ServiceListFailed = "CLAIM_LIST_CITY_V1_SERVICE_LIST_FAILED"
	MsgClaimListCityV1ServiceListFailed  = "Could not list city claims."

	CodeClaimPermV1HandlerParamsInvalid = "CLAIM_PERM_V1_HANDLER_PARAMS_INVALID"
	MsgClaimPermV1HandlerParamsInvalid  = "minecraftUuid, serverSlug, world, x, and z are required."

	CodeClaimPermV1ServiceFailed = "CLAIM_PERM_V1_SERVICE_FAILED"
	MsgClaimPermV1ServiceFailed  = "Could not check claim permission."
)

// Alliances
const (
	CodeAlliancePostV1HandlerBodyInvalid = "ALLIANCE_POST_V1_HANDLER_BODY_INVALID"
	MsgAlliancePostV1HandlerBodyInvalid  = "Request body must include leaderUuid, guildAId, and guildBId."

	CodeAlliancePostV1ServiceSameGuild = "ALLIANCE_POST_V1_SERVICE_SAME_GUILD"
	MsgAlliancePostV1ServiceSameGuild  = "Cannot ally a guild with itself."

	CodeAlliancePostV1ServiceLeaderNotFound = "ALLIANCE_POST_V1_SERVICE_LEADER_NOT_FOUND"
	MsgAlliancePostV1ServiceLeaderNotFound  = "Leader player not found."

	CodeAlliancePostV1ServiceGuildNotFound = "ALLIANCE_POST_V1_SERVICE_GUILD_NOT_FOUND"
	MsgAlliancePostV1ServiceGuildNotFound  = "Guild not found."

	CodeAlliancePostV1ServiceNotLeader = "ALLIANCE_POST_V1_SERVICE_NOT_LEADER"
	MsgAlliancePostV1ServiceNotLeader  = "Only guild A leader can propose alliance."

	CodeAlliancePostV1ServiceExists = "ALLIANCE_POST_V1_SERVICE_EXISTS"
	MsgAlliancePostV1ServiceExists  = "Alliance already exists between these guilds."

	CodeAlliancePostV1ServiceCreateFailed = "ALLIANCE_POST_V1_SERVICE_CREATE_FAILED"
	MsgAlliancePostV1ServiceCreateFailed  = "Could not create alliance."

	CodeAllianceListV1HandlerIDInvalid = "ALLIANCE_LIST_V1_HANDLER_ID_INVALID"
	MsgAllianceListV1HandlerIDInvalid  = "Guild id must be a valid UUID."

	CodeAllianceListV1ServiceGuildNotFound = "ALLIANCE_LIST_V1_SERVICE_GUILD_NOT_FOUND"
	MsgAllianceListV1ServiceGuildNotFound  = "Guild not found."

	CodeAllianceListV1ServiceListFailed = "ALLIANCE_LIST_V1_SERVICE_LIST_FAILED"
	MsgAllianceListV1ServiceListFailed  = "Could not list alliances."
)

// Audit
const (
	CodeAuditIngestV1HandlerBodyInvalid = "AUDIT_INGEST_V1_HANDLER_BODY_INVALID"
	MsgAuditIngestV1HandlerBodyInvalid  = "Request body must include a non-empty events array."

	CodeAuditIngestV1ServiceBatchTooLarge = "AUDIT_INGEST_V1_SERVICE_BATCH_TOO_LARGE"
	MsgAuditIngestV1ServiceBatchTooLarge  = "Audit batch cannot exceed 500 events."

	CodeAuditIngestV1ServiceFailed = "AUDIT_INGEST_V1_SERVICE_FAILED"
	MsgAuditIngestV1ServiceFailed  = "Could not ingest audit events."

	CodeAuditListV1HandlerIDInvalid = "AUDIT_LIST_V1_HANDLER_ID_INVALID"
	MsgAuditListV1HandlerIDInvalid  = "Player id must be a valid UUID."

	CodeAuditListV1ServicePlayerNotFound = "AUDIT_LIST_V1_SERVICE_PLAYER_NOT_FOUND"
	MsgAuditListV1ServicePlayerNotFound  = "Player not found."

	CodeAuditListV1ServiceListFailed = "AUDIT_LIST_V1_SERVICE_LIST_FAILED"
	MsgAuditListV1ServiceListFailed  = "Could not list audit events."
)

// Rollback
const (
	CodeRollbackPostV1HandlerBodyInvalid = "ROLLBACK_POST_V1_HANDLER_BODY_INVALID"
	MsgRollbackPostV1HandlerBodyInvalid  = "Request body must include targetUuid, actorUuid, serverSlug, and windowMinutes."

	CodeRollbackPostV1ServiceWindowInvalid = "ROLLBACK_POST_V1_SERVICE_WINDOW_INVALID"
	MsgRollbackPostV1ServiceWindowInvalid  = "windowMinutes must be between 1 and 1440."

	CodeRollbackPostV1ServiceTargetNotFound = "ROLLBACK_POST_V1_SERVICE_TARGET_NOT_FOUND"
	MsgRollbackPostV1ServiceTargetNotFound  = "Target player not found."

	CodeRollbackPostV1ServiceActorNotFound = "ROLLBACK_POST_V1_SERVICE_ACTOR_NOT_FOUND"
	MsgRollbackPostV1ServiceActorNotFound  = "Actor player not found."

	CodeRollbackPostV1ServiceCreateFailed = "ROLLBACK_POST_V1_SERVICE_CREATE_FAILED"
	MsgRollbackPostV1ServiceCreateFailed  = "Could not create rollback."

	CodeRollbackGetV1HandlerIDInvalid = "ROLLBACK_GET_V1_HANDLER_ID_INVALID"
	MsgRollbackGetV1HandlerIDInvalid  = "Rollback id must be a valid UUID."

	CodeRollbackGetV1ServiceNotFound = "ROLLBACK_GET_V1_SERVICE_NOT_FOUND"
	MsgRollbackGetV1ServiceNotFound  = "Rollback not found."

	CodeRollbackGetV1ServiceLoadFailed = "ROLLBACK_GET_V1_SERVICE_LOAD_FAILED"
	MsgRollbackGetV1ServiceLoadFailed  = "Could not load rollback."

	CodeRollbackItemsV1ServiceNotFound = "ROLLBACK_ITEMS_V1_SERVICE_NOT_FOUND"
	MsgRollbackItemsV1ServiceNotFound  = "Rollback not found."

	CodeRollbackItemsV1ServiceListFailed = "ROLLBACK_ITEMS_V1_SERVICE_LIST_FAILED"
	MsgRollbackItemsV1ServiceListFailed  = "Could not list rollback items."

	CodeRollbackCompleteV1HandlerBodyInvalid = "ROLLBACK_COMPLETE_V1_HANDLER_BODY_INVALID"
	MsgRollbackCompleteV1HandlerBodyInvalid  = "Request body must include appliedCount."

	CodeRollbackCompleteV1ServiceNotFound = "ROLLBACK_COMPLETE_V1_SERVICE_NOT_FOUND"
	MsgRollbackCompleteV1ServiceNotFound  = "Rollback not found."

	CodeRollbackCompleteV1ServiceFailed = "ROLLBACK_COMPLETE_V1_SERVICE_FAILED"
	MsgRollbackCompleteV1ServiceFailed  = "Could not complete rollback."
)

// Metrics
const (
	CodeMetricsOverviewV1ServiceFailed = "METRICS_OVERVIEW_V1_SERVICE_FAILED"
	MsgMetricsOverviewV1ServiceFailed  = "Could not load metrics overview."

	CodeMetricsTerritoryV1ServiceFailed = "METRICS_TERRITORY_V1_SERVICE_FAILED"
	MsgMetricsTerritoryV1ServiceFailed  = "Could not load territory metrics."
)

// Alerts
const (
	CodeAlertsListV1ServiceFailed = "ALERTS_LIST_V1_SERVICE_FAILED"
	MsgAlertsListV1ServiceFailed  = "Could not list alerts."

	CodeAlertsAckV1HandlerIDInvalid = "ALERTS_ACK_V1_HANDLER_ID_INVALID"
	MsgAlertsAckV1HandlerIDInvalid  = "Alert id must be a valid UUID."

	CodeAlertsAckV1ServiceNotFound = "ALERTS_ACK_V1_SERVICE_NOT_FOUND"
	MsgAlertsAckV1ServiceNotFound  = "Alert not found."

	CodeAlertsAckV1ServiceFailed = "ALERTS_ACK_V1_SERVICE_FAILED"
	MsgAlertsAckV1ServiceFailed  = "Could not acknowledge alert."
)

// Web link + session
const (
	CodeWebLinkV1HandlerBodyInvalid = "WEB_LINK_V1_HANDLER_BODY_INVALID"
	MsgWebLinkV1HandlerBodyInvalid  = "Request body must include playerId."

	CodeWebLinkV1ServiceFailed = "WEB_LINK_V1_SERVICE_FAILED"
	MsgWebLinkV1ServiceFailed  = "Could not create web link code."

	CodeWebSessionV1HandlerBodyInvalid = "WEB_SESSION_V1_HANDLER_BODY_INVALID"
	MsgWebSessionV1HandlerBodyInvalid  = "Request body must include code."

	CodeWebSessionV1HandlerCodeInvalid = "WEB_SESSION_V1_HANDLER_CODE_INVALID"
	MsgWebSessionV1HandlerCodeInvalid  = "Link code is invalid or expired."

	CodeWebSessionV1HandlerMissing = "WEB_SESSION_V1_HANDLER_MISSING"
	MsgWebSessionV1HandlerMissing  = "Session token is required."

	CodeWebSessionV1HandlerInvalid = "WEB_SESSION_V1_HANDLER_INVALID"
	MsgWebSessionV1HandlerInvalid  = "Session token is invalid or expired."
)

// Catalog + affiliation
const (
	CodeCatalogListV1ServiceFailed = "CATALOG_LIST_V1_SERVICE_FAILED"
	MsgCatalogListV1ServiceFailed  = "Could not load affiliation catalog."

	CodeAffiliationPatchV1HandlerBodyInvalid = "AFFILIATION_PATCH_V1_HANDLER_BODY_INVALID"
	MsgAffiliationPatchV1HandlerBodyInvalid  = "Request body is invalid."

	CodeAffiliationPatchV1ServiceFailed = "AFFILIATION_PATCH_V1_SERVICE_FAILED"
	MsgAffiliationPatchV1ServiceFailed  = "Could not update affiliation."

	CodeAffiliationPatchV1ServiceTypeInvalid = "AFFILIATION_PATCH_V1_SERVICE_TYPE_INVALID"
	MsgAffiliationPatchV1ServiceTypeInvalid  = "affiliationType must be student, staff, guest, or alumni."

	CodeAffiliationPatchV1ServiceCatalogIncomplete = "AFFILIATION_PATCH_V1_SERVICE_CATALOG_INCOMPLETE"
	MsgAffiliationPatchV1ServiceCatalogIncomplete  = "student and alumni require universitySlug, facultySlug, and courseSlug."

	CodeAffiliationPatchV1ServiceUniversityNotFound = "AFFILIATION_PATCH_V1_SERVICE_UNIVERSITY_NOT_FOUND"
	MsgAffiliationPatchV1ServiceUniversityNotFound  = "University not found in catalog."

	CodeAffiliationPatchV1ServiceFacultyNotFound = "AFFILIATION_PATCH_V1_SERVICE_FACULTY_NOT_FOUND"
	MsgAffiliationPatchV1ServiceFacultyNotFound  = "Faculty not found in catalog."

	CodeAffiliationPatchV1ServiceCourseNotFound = "AFFILIATION_PATCH_V1_SERVICE_COURSE_NOT_FOUND"
	MsgAffiliationPatchV1ServiceCourseNotFound  = "Course not found in catalog."

	CodeAffiliationPatchV1ServiceFacultyMismatch = "AFFILIATION_PATCH_V1_SERVICE_FACULTY_MISMATCH"
	MsgAffiliationPatchV1ServiceFacultyMismatch  = "Faculty does not belong to the selected university."

	CodeAffiliationPatchV1ServiceCourseMismatch = "AFFILIATION_PATCH_V1_SERVICE_COURSE_MISMATCH"
	MsgAffiliationPatchV1ServiceCourseMismatch  = "Course does not belong to the selected faculty."

	CodeAffiliationPatchV1ServiceGuestLocked = "AFFILIATION_PATCH_V1_SERVICE_GUEST_LOCKED"
	MsgAffiliationPatchV1ServiceGuestLocked  = "Guest accounts cannot change affiliation type."

	CodeGuildPostV1ServiceGuestLeader = "GUILD_POST_V1_SERVICE_GUEST_LEADER"
	MsgGuildPostV1ServiceGuestLeader  = "Guests cannot create or lead guilds."

	CodeCityPostV1ServiceGuestFounder = "CITY_POST_V1_SERVICE_GUEST_FOUNDER"
	MsgCityPostV1ServiceGuestFounder  = "Guests cannot create cities."

	CodeClaimPostV1ServiceGuest = "CLAIM_POST_V1_SERVICE_GUEST"
	MsgClaimPostV1ServiceGuest  = "Guests cannot create claims."
)
