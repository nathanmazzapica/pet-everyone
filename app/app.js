function logout() {
	alert("You have been logged out");
	localStorage.removeItem("token");
	window.location.href = "/";
}
