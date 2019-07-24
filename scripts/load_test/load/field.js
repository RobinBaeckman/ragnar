import http from "k6/http";
import { check, sleep } from "k6";

export let options = {
	vus: 300,
	duration: "10s"
	//stages: [
        //{ duration: "5m", target: 10 },
        //{ duration: "5m" }
        //{ duration: "5m", target: 0 },
    //]
};

export default function() {
	// Login as user
	var url = "http://127.0.0.1:3000/api/v1/fields/1d724594-5618-4e89-ad81-78bce5bedcba";
	var params =  { headers: { "Content-Type": "application/json" } }
	var res1 = http.get(url, params);
	check(res1, {
		"is status 200": (r) => r.status === 200
	});
	sleep(1);
};

