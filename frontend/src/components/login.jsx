import React, { useState, useRef, useContext } from 'react'
import { useRouter } from 'next/router'
import { LoginContext } from '@/components/Context'

export default function Login(props) {
  const [ errMsg, setErrMsg ] = useState('')
  const usernameRef = useRef('')
  const passwordRef = useRef('')
  const [ _, setInfo ] = useContext(LoginContext)
  const router = useRouter()

  const updateLogin = async (response) => {
    await setInfo({
      username: response.data.data.name,
      picture: response.data.picture
    })
  }

  const handleSubmit = (event) => { // function untuk menghandle onSubmit
    event.preventDefault() // untuk mengoperasikan submit melalui react (manual)
        
    var axios = require('axios'); //fetch POST API using axios
    var data = JSON.stringify({
      'username': usernameRef.current.value,
      'password': passwordRef.current.value
    })

    var config = {
      method: 'post',
      url: 'https://sirkelin.cloud/api/user/login/',
      headers: { 
        'Content-Type': 'application/json'
      },
      data: data
    }

    axios(config)
    .then(function (response) {
      if(response.status === axios.HttpStatusCode.Ok) {
        updateLogin(response)
          .then(function () {
            router.push("/messages")
          })
      }
    })
    .catch((error)=> {
      if(error?.response?.status === axios.HttpStatusCode.Unauthorized) {
        setErrMsg("Incorrect username and/or password")
      } else {
        setErrMsg("Server not responding")
      }
    })
  }

        
  return (
    <div className='mt-9 mx-auto'>
      <h1 className='text-white text-center font-extralight font-raleway'>Sirkelin</h1>
      <h2 className='pt-4 text-left text-xl mt-3 font-raleway'>Welcome Back... </h2>
      <p className={ errMsg === '' ? "offscreen" : "errMsg" }>{ errMsg }</p>
      <form onSubmit={ handleSubmit } className="mt-4">
        <div className="mb-2">
          <label
            htmlFor="email"
            className="block text-xs font-semibold"
          >
            Username
          </label>
          <input
            ref={usernameRef}
            id="username"
            autoComplete='off'
            type="text"
            className="block w-full px-4 py-1 mt-2 text-black bg-gray-dark rounded-md  focus:ring-purple-600 focus:outline-none focus:ring focus:ring-opacity-40"
          />
        </div>
        <div className="mb-2">
          <label
            htmlFor="password"
            className="block text-xs font-semibold mt-5"
          >
            Password
          </label>
          <input
            ref={passwordRef}
            id="password"
            type="password"
            autoComplete='off'
            className="block w-full px-4 py-1 mt-2 text-black bg-gray-dark rounded-md  focus:ring-purple-600 focus:outline-none focus:ring focus:ring-opacity-40"
          />
        </div>
        <div className="mt-6 text-center">
          <button className="inline-block px-6 py-2 border-2 border-purple-600 text-purple-600 font-medium text-xs leading-tight uppercase rounded-full hover:bg-gray-800 focus:outline-none focus:ring-0 transition duration-150 ease-in-out">
            Login
          </button>
        </div>
      </form>
        <p className="mt-8 text-xs font-light text-center text-gray-500">
        {" "}
        Don't have an account?{" "}
        <a
          onClick = { props.handler }
          className="font-medium text-purple-600 hover:underline cursor-pointer"
        >
          Sign up
        </a>
        </p>
    </div>
  )
}