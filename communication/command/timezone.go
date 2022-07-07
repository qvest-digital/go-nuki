package command

import (
	"fmt"
	"strconv"
)

type TimeZoneId uint16

var timeZoneIdMapping = map[TimeZoneId][]string{
	0:  {"Africa/Cairo", "UTC+2", "EET", "false"},
	1:  {"Africa/Lagos", "UTC+1", "WAT", "false"},
	2:  {"Africa/Maputo", "UTC+2", "CAT, SAST", "false"},
	3:  {"Africa/Nairobi", "UTC+3", "EAT", "false"},
	4:  {"America/Anchorage", "UTC-9/-8", "AKDT", "true"},
	5:  {"America/Argentina/Buenos_Aires", "UTC-3", "ART, UYT", "false"},
	6:  {"America/Chicago", "UTC-6/-5", "CDT", "true"},
	7:  {"America/Denver", "UTC-7/-6", "MDT", "true"},
	8:  {"America/Halifax", "UTC-4/-3", "ADT", "true"},
	9:  {"America/Los_Angeles", "UTC-8/-7", "PDT", "true"},
	10: {"America/Manaus", "UTC-4", "AMT, BOT, VET, AST, GYT", "false"},
	11: {"America/Mexico_City", "UTC-6/-5", "CDT", "true"},
	12: {"America/New_York", "UTC-5/-4", "EDT", "true"},
	13: {"America/Phoenix", "UTC-7", "MST", "false"},
	14: {"America/Regina", "UTC-6", "CST", "false"},
	15: {"America/Santiago", "UTC-4/-3", "CLST, AMST, WARST, PYST", "true"},
	16: {"America/Sao_Paulo", "UTC-3", "BRT", "false"},
	17: {"America/St_Johns", "UTC-3½/ -2½", "NDT", "true"},
	18: {"Asia/Bangkok", "UTC+7", "ICT, WIB", "false"},
	19: {"Asia/Dubai", "UTC+4", "SAMT, GET, AZT, GST, MUT, RET, SCT, AMT-Arm", "false"},
	20: {"Asia/Hong_Kong", "UTC+8", "HKT", "false"},
	21: {"Asia/Jerusalem", "UTC+2/+3", "IDT", "true"},
	22: {"Asia/Karachi", "UTC+5", "PKT, YEKT, TMT, UZT, TJT, ORAT", "false"},
	23: {"Asia/Kathmandu", "UTC+5¾", "NPT", "false"},
	24: {"Asia/Kolkata", "UTC+5½", "IST", "false"},
	25: {"Asia/Riyadh", "UTC+3", "AST-Arabia", "false"},
	26: {"Asia/Seoul", "UTC+9", "KST", "false"},
	27: {"Asia/Shanghai", "UTC+8", "CST, ULAT, IRKT, PHT, BND, WITA", "false"},
	28: {"Asia/Tehran", "UTC+3½", "ARST", "false"},
	29: {"Asia/Tokyo", "UTC+9", "JST, WIT, PWT, YAKT", "false"},
	30: {"Asia/Yangon", "UTC+6½", "MMT", "false"},
	31: {"Australia/Adelaide", "UTC+9½/10½", "ACDT", "true"},
	32: {"Australia/Brisbane", "UTC+10", "AEST, PGT, VLAT", "false"},
	33: {"Australia/Darwin", "UTC+9½", "ACST", "false"},
	34: {"Australia/Hobart", "UTC+10/+11", "AEDT", "true"},
	35: {"Australia/Perth", "UTC+8", "AWST", "false"},
	36: {"Australia/Sydney", "UTC+10/+11", "AEDT", "true"},
	37: {"Europe/Berlin", "UTC+1/+2", "CEST", "true"},
	38: {"Europe/Helsinki", "UTC+2/+3", "EEST", "true"},
	39: {"Europe/Istanbul", "UTC+3", "TRT", "false"},
	40: {"Europe/London", "UTC+0/+1", "BST, IST", "true"},
	41: {"Europe/Moscow", "UTC+3", "MSK", "false"},
	42: {"Pacific/Auckland", "UTC+12/+13", "NZDT", "true"},
	43: {"Pacific/Guam", "UTC+10", "ChST", "false"},
	44: {"Pacific/Honolulu", "UTC-10", "H(A)ST", "false"},
	45: {"Pacific/Pago_Pago", "UTC-11", "SST", "false"},
}

func (t TimeZoneId) Name() string {
	if entry, ok := timeZoneIdMapping[t]; ok {
		return entry[0]
	}
	return ""
}

func (t TimeZoneId) Offset() string {
	if entry, ok := timeZoneIdMapping[t]; ok {
		return entry[1]
	}
	return ""
}

func (t TimeZoneId) Timezone() string {
	if entry, ok := timeZoneIdMapping[t]; ok {
		return entry[2]
	}
	return ""
}

func (t TimeZoneId) DST() bool {
	if entry, ok := timeZoneIdMapping[t]; ok {
		b, _ := strconv.ParseBool(entry[3])
		return b
	}
	return false
}

func (t TimeZoneId) String() string {
	return fmt.Sprintf("%s | %s | %s | %v", t.Name(), t.Offset(), t.Timezone(), t.DST())
}
