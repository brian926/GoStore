import React from 'react'

function Signup() {
    const handleSubmit = (e) => {
        e.preventDefault()
        fetch('http://localhost:8000/users/signup', {
            method: 'POST',
            mode: 'no-cors',
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                first_name: `${e.target.first_name.value}`,
                last_name: `${e.target.last_name.value}`,
                password: `${e.target.password.value}`,
                email: `${e.target.email.value}`,
                phone: `${e.target.phone.value}`,
            }),
        })
        .then((res) => {
            console.log(res)
        })
        .catch((err) => {
            console.log(err)
        })
    }

    return (
        <div>
            <div >
                <form onSubmit={handleSubmit}>
                        <input
                            type='text'
                            name='first_name'
                            placeholder='first_name'
                        />
                        <input
                            type='text'
                            name='last_name'
                            placeholder='last_name'
                        />
                        <input 
                            type='text'
                            name='password'
                            placeholder='password'
                        />
                        <input 
                            type='text'
                            name='email'
                            placeholder='email'
                        />
                        <input 
                            type='text'
                            name='phone'
                            placeholder='phone'
                        />
                        <input type='submit' value='submit' />
                </form>
            </div>
        </div>
    )
}

export default Signup