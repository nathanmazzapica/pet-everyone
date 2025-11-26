document.addEventListener("DOMContentLoaded", (event) => {
	if (!localStorage.getItem("token")) {
		alert("You must be logged in to create a pet");
		window.location.href = "/login";
	}
});

const createForm = document.getElementById('new-pet-form')

createForm.addEventListener('submit', e => {
    e.preventDefault()
    uploadPet()
})

async function uploadPet() {
    const petName = document.getElementById('pet-name-input').value;
    const petImageFile = document.getElementById('pet-image-input').files[0];

    if (!petImageFile) {
        return;
    }

    const formData = new FormData();
    formData.append('petName', petName);
    formData.append('petImageFile', petImageFile);

    try {
        const res = await fetch(`/api/create`, {
            method: 'POST',
			headers: {
				Authorization: `Bearer ${localStorage.getItem('token')}`,
			},
            body: formData,
        });
        if (!res.ok) {
            const data = await res.json();
            throw new Error(`Failed to create pet. Error: ${data.error}`);
        }
        const data = await res.json();
        console.table(data)
        console.log('Successfully created');
        window.location.href = `/pet/${data.PetID}`
    } catch(error) {
        alert(`Failed to create pet. Error: ${error}`);
    }
}