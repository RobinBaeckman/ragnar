import http from "k6/http";
import { check, sleep } from "k6";

export let options = {
	vus: 200,
	duration: "10s"
	//stages: [
        //{ duration: "5m", target: 10 },
        //{ duration: "5m" }
        //{ duration: "5m", target: 0 },
    //]
};

export default function() {
	// Login as user
	var url = "http://127.0.0.1:3000/api/v1/sessions";
	var payload = JSON.stringify({ email: "aa_static_user@mail.com", password: "Secret123" });
	var params =  { headers: { "Content-Type": "application/json" } }
	var res1 = http.post(url, payload, params);

	// Show Messages 
	var url2 = "http://127.0.0.1:3000/api/v2/messages?messagable_id=20d8d32f-0ab3-4e31-9b71-13496dea4e19";
	var params2 =  { headers: { "Content-Type": "application/json", "X-Farmfes-Session-Id": String(res1.headers["X-Farmfes-Session-Id"])} }
	
	var res2 = http.get(url2, params2);
	check(res2, {
		"is status 200": (r) => r.status === 200
	});

	var obj = JSON.parse(String(res2.body))
	//console.log(obj[0].body);
	sleep(2);
};

