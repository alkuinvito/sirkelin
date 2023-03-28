import React, { useState } from "react"

export default function GetUserModal() {
  const [ result, setResult ] = useState([])

  const getUsers = (query) => {
    const axios = require('axios')
    axios.get(process.env.NEXT_PUBLIC_APP_HOST + '/api/user/search?q=' + query, { withCredentials: true })
    .then(response => {
      if (response.data.data.users === null) {
        setResult([])
      } else {
        setResult(response.data.data.users)
      }
    })
    .catch(console.error)
  }

  let delay = 0
  const handleChange = (e) => {
    if (e.target.value.length >= 3 && e.target.value.length <= 16) {
      clearTimeout(delay)
      delay = setTimeout(() => {
        getUsers(e.target.value)
      }, 320)
    } else {
      setResult([])
    }
  }

  return (
    <div className='absolute w-48 h-32 p-4 rounded-lg bg-black'>
      <input
        type='text'
        name='fullname'
        className='grow h-7 text-sm bg-transparent text-gray-200 focus:outline-none'
        placeholder='Search contacts...'
        autoComplete="off"
        onChange={handleChange}
        />
        <ul>
          { result.map((data)=><li key={data.ID} className='flex gap-2'><img className='rounded-full' src={data.Picture} width={24} height={24}/>{data.Fullname}</li>) }
        </ul>
    </div>
  )
}