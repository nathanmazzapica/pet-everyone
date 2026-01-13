function logout() {
async function logout() {
	await fetch("/api/logout", {
		method: "POST",
		credentials: "include",
	});
	alert("You have been logged out");
	window.location.href = "/";
}

function hasSession() {
	return document.cookie.split(";").some((c) => c.trim().startsWith("session_present="));
}

function wireAuthLink() {
	const authLink = document.getElementById("auth-link");
	if (!authLink) return;

	if (hasSession()) {
		authLink.textContent = "Logout";
		authLink.href = "#";
		authLink.addEventListener("click", async (e) => {
			e.preventDefault();
			await logout();
		});
	} else {
		authLink.textContent = "Login";
		authLink.href = "/login";
	}
}

document.addEventListener("DOMContentLoaded", () => {
	wireAuthLink();
});
