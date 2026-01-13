(() => {
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

	async function fetchDisplayName() {
		try {
			const res = await fetch("/api/display-name", {
				method: "GET",
				credentials: "include",
			});
			if (!res.ok) return null;
			const data = await res.json();
			if (data && typeof data.display_name === "string") {
				return data.display_name;
			}
			return null;
		} catch (err) {
			console.warn("Failed to fetch display name", err);
			return null;
		}
	}

	async function setWelcomeHeader() {
		const header = document.getElementById("welcome-header");
		if (!header) return;
		const name = await fetchDisplayName();
		header.textContent = name ? `Welcome, ${name}` : "Welcome";
	}

	document.addEventListener("DOMContentLoaded", () => {
		wireAuthLink();
		setWelcomeHeader();
	});
})();
