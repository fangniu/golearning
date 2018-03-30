package main

import (
	"fmt"
	"regexp"
	//"strings"
	"strings"
	"errors"
)

var (
	noServiceNameFieldError = errors.New("没有service name字段")
	measurementRe *regexp.Regexp
	serviceNameRe *regexp.Regexp
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

func checkQuery(query string) error {
	if measurementRe == nil {
		serviceNameRe = regexp.MustCompile(fmt.Sprintf(`("%v" =~ /\^([^\$]+)\$/)|("%v" = '([^ ]+)')`, "service_name", "service_name"))
		var measurements []string
		ss := []string{`drop_data"."data`}
		for _, s := range ss {
			measurements = append(measurements, fmt.Sprintf(`("%v")`, s))
		}
		measurementRe = regexp.MustCompile(".*FROM " + strings.Join(measurements, "|"))
		//fmt.Println(".*FROM " + strings.Join(measurements, "|"))
	}
	if !measurementRe.MatchString(query) {
		return nil
	}

	//fmt.Println(fmt.Sprintf(`"%v" =~ /\^([0-9a-zA-Z_]+)\$/`, "service_name"))
	if len(serviceNameRe.FindAllStringSubmatch(query, -1)) == 0 {
		return noServiceNameFieldError
	}


	for _, q := range strings.Split(query, ";") {
		if !measurementRe.MatchString(q) {
			continue
		}
		serviceNames := serviceNameRe.FindAllStringSubmatch(q, -1)
		if len(serviceNames) == 0 {
			return noServiceNameFieldError
		}
		for _, arr := range serviceNames {
			if len(arr[1]) == 0 {
				fmt.Println(strings.Replace(arr[4], "\\", "", -1))
			} else {
				fmt.Println(strings.Replace(arr[2], "\\", "", -1))
			}
		}

	}
	return nil
}

func validServiceName(name string, serviceNames []string) bool {
	for _, n := range serviceNames {
		if name == n {
			return true
		}
	}
	return false
}



func main() {
	//s1 := `SELECT "value" FROM "drop_data"."data" WHERE ("service_name" = 'serviceAB' AND "attr_name" =~ /^NTL\.Buffer\.Use$/ AND "service_name" =~ /^valet_property_equip_preview$/  AND "type" = '1') AND time >= now() - 5m`
	//s1 := `SELECT "value" FROM "drop_data"."data" WHERE ("service_name" = 'serviceADCa') AND time >= now() - 15m`
	//s1 := `SELECT "value" FROM "drop_data"."data" WHERE ("service_name" =~ /^valet_property_equip_preview$/) AND time >= now() - 15m`
	s1 := `SELECT "value" FROM "drop_data"."data" WHERE ("service_name" = 'aaa' AND "attr_name" =~ /^Dealing\.PerSecond$/ AND "server_id" =~ /^0$/ AND "type" = '1') AND time >= now() - 5m;SELECT "value" FROM "drop_data"."data" WHERE ("service_name" =~ /^DataRouteSystem\.GroupRoute$/ AND "attr_name" =~ /^Dealing\.PerSecond$/ AND "server_id" =~ /^0$/ AND "type" = '2') AND time >= now() - 5m;SELECT "value" FROM "drop_data"."data" WHERE ("service_name" =~ /^DataRouteSystem\.GroupRoute$/ AND "attr_name" =~ /^Dealing\.PerSecond$/ AND "server_id" =~ /^0$/ AND "type" = '3') AND time >= now() - 5m;SELECT "value" FROM "drop_data"."data" WHERE ("service_name" =~ /^DataRouteSystem\.GroupRoute$/ AND "attr_name" =~ /^Dealing\.PerSecond$/ AND "server_id" =~ /^0$/ AND "type" = '4') AND time >= now() - 5m;SELECT "value" FROM "drop_data"."data" WHERE ("service_name" =~ /^DataRouteSystem\.GroupRoute$/ AND "attr_name" =~ /^Dealing\.PerSecond$/ AND "server_id" =~ /^0$/ AND "type" = '5') AND time >= now() - 5m;SELECT "value" FROM "drop_data"."data" WHERE ("service_name" =~ /^DataRouteSystem\.GroupRoute$/ AND "attr_name" =~ /^Dealing\.PerSecond$/ AND "server_id" =~ /^0$/ AND "type" = '6') AND time >= now() - 5m;SELECT "value" FROM "drop_data"."data" WHERE ("service_name" =~ /^DataRouteSystem\.GroupRoute$/ AND "attr_name" =~ /^Dealing\.PerSecond$/ AND "server_id" =~ /^0$/ AND "type" = '7') AND time >= now() - 5m;SELECT "value" AS "shift_1_days" FROM "drop_data"."data" WHERE ("service_name" =~ /^DataRouteSystem\.GroupRoute$/ AND "attr_name" =~ /^Dealing\.PerSecond$/ AND "server_id" =~ /^0$/ AND "type" = '1') AND time >= now() - 5m;SELECT "value" AS "shift_1_days" FROM "drop_data"."data" WHERE ("service_name" =~ /^DataRouteSystem\.GroupRoute$/ AND "attr_name" =~ /^Dealing\.PerSecond$/ AND "server_id" =~ /^0$/ AND "type" = '2') AND time >= now() - 5m;SELECT "value" AS "shift_1_days" FROM "drop_data"."data" WHERE ("service_name" =~ /^DataRouteSystem\.GroupRoute$/ AND "attr_name" =~ /^Dealing\.PerSecond$/ AND "server_id" =~ /^0$/ AND "type" = '3') AND time >= now() - 5m;SELECT "value" AS "shift_1_days" FROM "drop_data"."data" WHERE ("service_name" =~ /^DataRouteSystem\.GroupRoute$/ AND "attr_name" =~ /^Dealing\.PerSecond$/ AND "server_id" =~ /^0$/ AND "type" = '4') AND time >= now() - 5m;SELECT "value" AS "shift_1_days" FROM "drop_data"."data" WHERE ("service_name" =~ /^DataRouteSystem\.GroupRoute$/ AND "attr_name" =~ /^Dealing\.PerSecond$/ AND "server_id" =~ /^0$/ AND "type" = '5') AND time >= now() - 5m;SELECT "value" AS "shift_1_days" FROM "drop_data"."data" WHERE ("service_name" =~ /^DataRouteSystem\.GroupRoute$/ AND "attr_name" =~ /^Dealing\.PerSecond$/ AND "server_id" =~ /^0$/ AND "type" = '6') AND time >= now() - 5m;SELECT "value" AS "shift_1_days" FROM "drop_data"."data" WHERE ("service_name" =~ /^DataRouteSystem\.GroupRoute$/ AND "attr_name" =~ /^Dealing\.PerSecond$/ AND "server_id" =~ /^0$/ AND "type" = '7') AND time >= now() - 5m`
	err := checkQuery(s1)
	fmt.Println(err)
}

