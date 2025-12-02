package audit

// System is an application system targeted for audition.
type System string

const (
	SystemBadges                       System = "badges"
	SystemBadgeUsers                   System = "badge_users"
	SystemChannelsCommands             System = "channels_commands"
	SystemChannelsCommandGroups        System = "channels_command_groups"
	SystemChannelsCustomVars           System = "channels_custom_vars"
	SystemChannelsGames8Ball           System = "channels_games_8ball"
	SystemChannelsGamesDuel            System = "channels_games_duel"
	SystemChannelsGamesRussianRoulette System = "channels_games_russian_roulette"
	SystemChannelsGamesSeppuku         System = "channels_games_seppuku"
	SystemChannelsGamesVoteban         System = "channels_games_voteban"
	SystemChannelsGreetings            System = "channels_greetings"
	SystemChannelsKeywords             System = "channels_keywords"
	SystemChannelsModerationSettings   System = "channels_moderation_settings"
	SystemChannelsOverlaysChat         System = "channels_overlays_chat"
	SystemChannelsOverlaysDudes        System = "channels_overlays_dudes"
	SystemChannelsOverlaysKappaGen     System = "channels_overlays_kappa_gen"
	SystemChannelsOverlaysNowPlaying   System = "channels_overlays_now_playing"
	SystemChannelsRoles                System = "channels_roles"
	SystemChannelsTimers               System = "channels_timers"
	SystemChannelsSongRequestsSettings System = "channels_song_requests_settings"
	SystemChannelsIntegrations         System = "channels_integrations"
	SystemChannelsAlerts               System = "channels_alerts"
	SystemChannelsChatAlerts           System = "channels_chat_alerts"
	SystemChannelsChatTranslation      System = "channels_chat_translation"
)
