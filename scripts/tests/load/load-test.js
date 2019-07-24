import http from "k6/http";
import { check, sleep } from "k6";
import { Rate, Trend } from "k6/metrics";

let createUserErrorRate = new Rate("Create User errors");
let loginErrorRate = new Rate("Login errors");
let readUserErrorRate = new Rate("Read User errors");
let readUsersErrorRate = new Rate("Read Users errors");
let updateUserErrorRate = new Rate("Update User errors");
let deleteUserErrorRate = new Rate("Delete User errors");

let createUserTrend = new Trend("Create User");
let loginTrend = new Trend("Login");
let readUserTrend = new Trend("Read User");
let readUsersTrend = new Trend("Read Users");
let updateUserTrend = new Trend("Update User");
let deleteUserTrend = new Trend("Delete User");

export let options = {
	vus: 40,
	duration: "60s",
	thresholds: {
    "Create User": ["p(95)<500"],
    "Login": ["p(95)<800"],
    "Read User": ["p(95)<500"],
    "Read Users": ["p(95)<500"],
    "Update User": ["p(95)<800"],
    "Delete User": ["p(95)<500"],
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

	var createUserResp = http.post(urlUsers, createUserData, defaultParams);
	var loginResp = http.post(urlLogin, loginData, defaultParams);

	check(createUserResp, {
    "status is 201": (r) => r.status == 201
  	}) || createUserErrorRate.add(1);  
  	createUserTrend.add(createUserResp.timings.duration);

	check(loginResp, {
    "status is 200": (r) => r.status == 200
  	}) || loginErrorRate.add(1);  
  	loginTrend.add(loginResp.timings.duration);

	// Post Login
	var cookie = String(loginResp.headers["cookie"])
	var user = JSON.parse(createUserResp.body)
	var urlUser = urlUsers+"/"+user.id;
	var createUserData = JSON.stringify({ email: email, password: password, firstName: "UpdatedRolf", lastName: "UpdatedBaeckman" });
	var authParams = { headers: { "Content-Type": "application/json", "cookie": cookie} }

	var readUserResp = http.get(urlUser, authParams);
	var readUsersResp = http.get(urlUsers, authParams);
	var updateUserResp = http.put(urlUsers, createUserData, defaultParams);
	var deleteUserResp = http.del(urlUser, authParams);

	check(readUserResp, {
    "status is 200": (r) => r.status == 200
  	}) || readUserErrorRate.add(1);  
  	readUserTrend.add(readUserResp.timings.duration);

	check(updateUserResp, {
    "status is 200": (r) => r.status == 200
  	}) || updateUserErrorRate.add(1);  
  	updateUserTrend.add(updateUserResp.timings.duration);

	check(deleteUserResp, {
    "status is 200": (r) => r.status == 200
  	}) || deleteUserErrorRate.add(1);  
  	deleteUserTrend.add(deleteUserResp.timings.duration);

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
