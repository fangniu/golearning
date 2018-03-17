package main

import (
	"fmt"
	"regexp"
	//"strings"
)

func f1()  {
	// Compile the expression once, usually at init time.
	// Use raw strings to avoid having to quote the backslashes.
	//var validID = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`)
	//
	//fmt.Println(validID.MatchString("adam[23]"))
	//fmt.Println(validID.MatchString("eve[7]"))
	//fmt.Println(validID.MatchString("Job[48]"))
	//fmt.Println(validID.MatchString("snakey"))

	var validMeasurement = regexp.MustCompile(`.*FROM ("service_retention1"."service")|("service_retention"."1service")`)
	var re = regexp.MustCompile(`"service_name" =~ /\^([0-9a-zA-Z_]+)\$/`)
	s1 := `SELECT "value" FROM "service_retention"."service" WHERE ("service_name" =~ /^valet_property_equip_preview$/ AND "attr_name" =~ /^dealValetCacheClearRQ\.RQCOUNT$/ AND "service_id" =~ /^1024119$/ AND "repotor_ip" =~ /^10\.33\.109\.104$/ AND "op_type" = '1') AND time >= now() - 1h;SELECT "value" FROM "service_retention"."service" WHERE ("service_name" =~ /^valet_property_equip_preview$/ AND "attr_name" =~ /^dealValetCacheClearRQ\.RQCOUNT$/ AND "service_id" =~ /^1024119$/ AND "repotor_ip" =~ /^10\.33\.109\.104$/ AND "op_type" = '2') AND time >= now() - 1h;SELECT "value" FROM "service_retention"."service" WHERE ("service_name" =~ /^valet_property_equip_preview$/ AND "attr_name" =~ /^dealValetCacheClearRQ\.RQCOUNT$/ AND "service_id" =~ /^1024119$/ AND "repotor_ip" =~ /^10\.33\.109\.104$/ AND "op_type" = '3') AND time >= now() - 1h;SELECT "value" FROM "service_retention"."service" WHERE ("service_name" =~ /^valet_property_equip_preview$/ AND "attr_name" =~ /^dealValetCacheClearRQ\.RQCOUNT$/ AND "service_id" =~ /^1024119$/ AND "repotor_ip" =~ /^10\.33\.109\.104$/ AND "op_type" = '4') AND time >= now() - 1h;SELECT "value" FROM "service_retention"."service" WHERE ("service_name" =~ /^valet_property_equip_preview$/ AND "attr_name" =~ /^dealValetCacheClearRQ\.RQCOUNT$/ AND "service_id" =~ /^1024119$/ AND "repotor_ip" =~ /^10\.33\.109\.104$/ AND "op_type" = '5') AND time >= now() - 1h;SELECT "value" FROM "service_retention"."service" WHERE ("service_name" =~ /^valet_property_equip_preview$/ AND "attr_name" =~ /^dealValetCacheClearRQ\.RQCOUNT$/ AND "service_id" =~ /^1024119$/ AND "repotor_ip" =~ /^10\.33\.109\.104$/ AND "op_type" = '6') AND time >= now() - 1h;SELECT "value" FROM "service_retention"."service" WHERE ("service_name" =~ /^valet_property_equip_preview$/ AND "attr_name" =~ /^dealValetCacheClearRQ\.RQCOUNT$/ AND "service_id" =~ /^1024119$/ AND "repotor_ip" =~ /^10\.33\.109\.104$/ AND "op_type" = '7') AND time >= now() - 1h`
	var services = map[string]int{}
	if validMeasurement.MatchString(s1) {
		for _, arr := range re.FindAllStringSubmatch(s1, -1) {
			fmt.Println(arr[1])
			services[arr[1]] = 0
		}
	}

	//fmt.Println(re.FindAllStringSubmatch(s1, -1))
}

func f2()  {
	var re *regexp.Regexp
	fmt.Println(re)
}

func main() {

	f1()
}

