import http from "k6/http";
import { check, sleep } from "k6";

export let options = {
	vus: 10,
	duration: "3m",
	//iterations: 1
	//stages: [
        //{ duration: "5m", target: 10 },
        //{ duration: "5m" }
        //{ duration: "5m", target: 0 },
    //]
};

export default function() {
	// Login as user
	var sessionUrl = "http://127.0.0.1:3000/api/v1/sessions";
	var sessionPayload = JSON.stringify({ email: "aa_static_user@mail.com", password: "Secret123" });
	var sessionParams =  { headers: { "Content-Type": "application/json" } }
	var sessionRes = http.post(sessionUrl, sessionPayload, sessionParams);
	
	// Show messages 
	var messageUrl = "http://127.0.0.1:3000/api/v2/messages?messagable_id=20d8d32f-0ab3-4e31-9b71-13496dea4e19";
	var messageParams =  { headers: { "Content-Type": "application/json", "X-Farmfes-Session-Id": String(sessionRes.headers["X-Farmfes-Session-Id"])} }
	var messageRes = http.get(messageUrl, messageParams);
	check(messageRes, {
		"is status 200": (r) => r.status === 200
	});

	//var obj = JSON.parse(String(messageRes.body))
	//console.log(obj[0].body);
	sleep(1);
};

