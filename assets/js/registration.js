let SubmitButton = document.querySelector('.form')
if (SubmitButton) {
    SubmitButton.addEventListener('click', () => {
        fetch("http://127.0.0.1:8080/api/register", {
            method: "POST",
            headers: {
                "Accept": "application/json",
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                "username": document.getElementById('username').value,
                "password": document.getElementById('password').value
            })
        })
            .then(response => response.json())
            .then(user => console.log(user))
    })
}