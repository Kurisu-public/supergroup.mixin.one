package config

import (
	"io/ioutil"
	"log"
	"path"

	"gopkg.in/yaml.v2"
)

const ConfigFile = "config.yaml"
const BuildVersion = "BUILD_VERSION"
const DefaultColor = "#46B8DA"

type PaymentAsset struct {
	Symbol  string `yaml:"symbol" json:"symbol"`
	AssetId string `yaml:"asset_id" json:"asset_id"`
	Amount  string `yaml:"amount" json:"amount"`
}

type Shortcut struct {
	Icon    string `yaml:"icon" json:"icon"`
	LabelEn string `yaml:"label_en" json:"label_en"`
	LabelZh string `yaml:"label_zh" json:"label_zh"`
	Url     string `yaml:"url" json:"url"`
}

type ShortcutGroup struct {
	LabelEn string     `yaml:"label_en" json:"label_en"`
	LabelZh string     `yaml:"label_zh" json:"label_zh"`
	Items   []Shortcut `yaml:"shortcuts" json:"shortcuts"`
}

type Config struct {
	Service struct {
		Name             string `yaml:"name"`
		Environment      string `yaml:"enviroment"`
		HTTPListenPort   int    `yaml:"port"`
		HTTPResourceHost string `yaml:"host"`
	} `yaml:"service"`
	Database struct {
		User     string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"database_name"`
	} `yaml:"database"`
	System struct {
		MessageShardModifier                       string   `yaml:"message_shard_modifier"`
		MessageShardSize                           int64    `yaml:"message_shard_size"`
		PriceAssetsEnable                          bool     `yaml:"price_asset_enable"`
		MinimumUsdtPrice                           string   `yaml:"minimum_usdt_price"`
		MaximumPacketNumber                        int64    `yaml:"maximum_packet_number"`
		AudioMessageEnable                         bool     `yaml:"audio_message_enable"`
		ImageMessageEnable                         bool     `yaml:"image_message_enable"`
		VideoMessageEnable                         bool     `yaml:"video_message_enable"`
		ContactMessageEnable                       bool     `yaml:"contact_message_enable"`
		LimitMessageDuration                       int64    `yaml:"limit_message_duration"`
		LimitMessageNumber                         int      `yaml:"limit_message_number"`
		DetectQRCodeEnabled                        bool     `yaml:"detect_image"`
		DetectLinkEnabled                          bool     `yaml:"detect_link"`
		KeywordReplyEnable                         bool     `yaml:"keyword_reply_enable"`
		ImmediateDeleteExpiredDistributedMsgEnable bool     `yaml:"immediate_delete_expired_distributed_msg_enable"`
		SuperOperatorList                          []string `yaml:"super_operator_list"`
		SuperOperators                             map[string]bool
		OperatorList                               []string `yaml:"operator_list"`
		Operators                                  map[string]bool
		PayToJoin                                  bool           `yaml:"pay_to_join"`
		AccpetPaymentAssetList                     []PaymentAsset `yaml:"accept_asset_list"`
	} `yaml:"system"`
	Appearance struct {
		HomeWelcomeMessage string          `yaml:"home_welcome_message"`
		HomeShortcutGroups []ShortcutGroup `yaml:"home_shortcut_groups"`
	} `yaml:"appearance"`
	MessageTemplate struct {
		WelcomeMessage          string         `yaml:"welcome_message"`
		MessageTipsGuest        string         `yaml:"message_tips_guest"`
		MessageTipsHelp         string         `yaml:"message_tips_help"`
		GroupRedPacket          string         `yaml:"group_redpacket"`
		GroupRedPacketShortDesc string         `yaml:"group_redpacket_short_desc"`
		GroupRedPacketDesc      string         `yaml:"group_redpacket_desc"`
		GroupOpenedRedPacket    string         `yaml:"group_opened_redpacket"`
		MessageProhibit         string         `yaml:"message_prohibit"`
		MessageAllow            string         `yaml:"message_allow"`
		MessageTipsJoin         string         `yaml:"message_tips_join"`
		MessageTipsHelpBtn      string         `yaml:"message_tips_help_btn"`
		MessageTipsUnsubscribe  string         `yaml:"message_tips_unsubscribe"`
		MessageRewardLabel      string         `yaml:"message_reward_label"`
		MessageRewardMemo       string         `yaml:"message_reward_memo"`
		MessageTipsTooMany      string         `yaml:"message_tips_too_many"`
		MessageCommandsInfo     string         `yaml:"message_commands_info"`
		MessageCommandsInfoResp string         `yaml:"message_commands_info_resp"`
		KeywordReplyList        []KeywordReply `yaml:"keyword_reply_list"`
		Keywords                map[string][]KeywordReplyMessage
	} `yaml:"message_template"`
	Mixin struct {
		ClientId        string `yaml:"client_id"`
		ClientSecret    string `yaml:"client_secret"`
		SessionAssetPIN string `yaml:"session_asset_pin"`
		PinToken        string `yaml:"pin_token"`
		SessionId       string `yaml:"session_id"`
		SessionKey      string `yaml:"session_key"`
	} `yaml:"mixin"`
}

type ExportedConfig struct {
	MixinClientId          string          `json:"mixin_client_id"`
	HTTPResourceHost       string          `json:"host"`
	AccpetPaymentAssetList []PaymentAsset  `json:"accept_asset_list"`
	HomeWelcomeMessage     string          `json:"home_welcome_message"`
	HomeShortcutGroups     []ShortcutGroup `json:"home_shortcut_groups"`
}

type KeywordReply struct {
	Keyword  string                `yaml:"keyword"`
	Messages []KeywordReplyMessage `yaml:"messages"`
}

type KeywordReplyMessage struct {
	Category string `yaml:"category"`
	Data     string `yaml:"data"`
}

type AppButtonGroup struct {
	Label  string `yaml:"label"`
	Color  string `yaml:"color"`
	Action string `yaml:"action"`
}

type AppCard struct {
	IconURL     string `yaml:"icon_url"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Action      string `yaml:"action"`
}

var AppConfig *Config

func LoadConfig(dir string) {
	data, err := ioutil.ReadFile(path.Join(dir, ConfigFile))
	if err != nil {
		log.Panicln(err)
	}
	AppConfig = &Config{}
	err = yaml.Unmarshal(data, AppConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// operators
	AppConfig.System.Operators = make(map[string]bool)
	for _, op := range AppConfig.System.OperatorList {
		AppConfig.System.Operators[op] = true
	}
	// super operators
	AppConfig.System.SuperOperators = make(map[string]bool)
	for _, sop := range AppConfig.System.SuperOperatorList {
		AppConfig.System.SuperOperators[sop] = true
	}
	// keywords
	AppConfig.MessageTemplate.Keywords = make(map[string][]KeywordReplyMessage)
	for _, kw := range AppConfig.MessageTemplate.KeywordReplyList {
		AppConfig.MessageTemplate.Keywords[kw.Keyword] = kw.Messages
	}
}

func GetExported() ExportedConfig {
	return ExportedConfig{
		MixinClientId:          AppConfig.Mixin.ClientId,
		HTTPResourceHost:       AppConfig.Service.HTTPResourceHost,
		AccpetPaymentAssetList: AppConfig.System.AccpetPaymentAssetList,
		HomeWelcomeMessage:     AppConfig.Appearance.HomeWelcomeMessage,
		HomeShortcutGroups:     AppConfig.Appearance.HomeShortcutGroups,
	}
}
