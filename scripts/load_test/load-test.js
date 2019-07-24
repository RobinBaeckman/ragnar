import http from "k6/http";
import { check, sleep } from "k6";
import { Rate, Trend } from "k6/metrics";

let createUserErrorRate = new Rate("Create User errors");
let loginErrorRate = new Rate("Login errors");
let readUserErrorRate = new Rate("Read User errors");
let createUserTrend = new Trend("Create User");
let loginTrend = new Trend("Login");
let readUserTrend = new Trend("Read User");

export let errorRate = new Rate("errors");

export let options = {
	vus: 20,
	duration: "30s",
	thresholds: {
    "Create User": ["p(95)<500"],
    "Login": ["p(95)<800"],
    "Read User": ["p(95)<500"],
  }
};

export default function() {
	// Pre Login
	var urlUsers = "http://127.0.0.1:3000/v1/users";
	var urlLogin = "http://127.0.0.1:3000/v1/login";
	//const email = `user+${__VU}@mail.com`;
	var email = generateRandomString(10)+"@mail.com";
	const password = "secret"
	var createUserData = JSON.stringify({ email: email, password: password, firstName: "Rolf", lastName: "Baeckman" });
	var loginData = JSON.stringify({ email: email, password: password });
	var defaultParams =  { headers: { "Content-Type": "application/json" } }

	//var createUserResp = http.post(urlUsers, createUserData, defaultParams);

	let requests = {
    "Crete User": {
      method: "POST",
      url: urlUsers,
      params: defaultParams,
      body: createUserData
    },
    "Login": {
      method: "POST",
      url: urlLogin,
      params: defaultParams,
      body: loginData
    },
  	};
	let responses = http.batch(requests);
  	let createUserResp = responses["Create User"];
	let loginResp = responses["Login"];

	check(createUserResp, {
    "status is 201": (r) => r.status == 201
  	}) || createUserErrorRate.add(1);  
  	createUserTrend.add(createUserResp.timings.duration);

	check(loginResp, {
    "status is 200": (r) => r.status == 200
  	}) || loginErrorRate.add(1);  
  	loginTrend.add(loginResp.timings.duration);

	// Post Login
	var user = JSON.parse(String(createUserResp.body))
	var urlReadUser = urlUsers+"/"+user.id;
	//console.log(urlReadUser);
	var authParams =  { headers: { "Content-Type": "application/json", "cookie": String(loginRes.headers["cookie"])} }
	let requests2 = {
    "Read User": {
      method: "GET",
      url: urlReadUser,
      params: authParams
    },
  	};

	let responses2 = http.batch(requests2);
	let readUserResp = responses2["Read User"];

	check(readUserResp, {
    "status is 200": (r) => r.status == 200
  	}) || readUserErrorRate.add(1);  
  	readUserTrend.add(readUserResp.timings.duration);

	sleep(1);
};

function generateRandomString(length) {
   var result           = '';
   var characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
   var charactersLength = characters.length;
   for ( var i = 0; i < length; i++ ) {
      result += characters.charAt(Math.floor(Math.random() * charactersLength));
   }
   return result;
}
