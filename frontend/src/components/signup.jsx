import React, { useState, useRef } from 'react'

export default function Signup(props) {
  const [ errMsg, setErrMsg ] = useState('')
  const usernameRef = useRef('')
  const emailRef = useRef('')
  const passwordRef = useRef('')
  function handleSubmit(event) { // function untuk menghandle onSubmit
    event.preventDefault() // untuk mengoperasikan submit melalui react (manual)
    var axios = require('axios'); //fetch POST API using axios
    var data = JSON.stringify({
      'email': emailRef.current.value,
      'username': usernameRef.current.value,
      'password': passwordRef.current.value
    });
    var config = {
      method: 'post',
      url: 'https://malakh.space/api/user/register/',
      headers: {
        'Content-Type': 'application/json'
      },
      data : data
    };
    axios(config)
    .then(function (response) {
       console.log(JSON.stringify(response.data));
    })
    .catch((error)=> {
      if(error.response.status === axios.HttpStatusCode.BadRequest) {
        setErrMsg(error.response.data.data.error)
      }
    });
  }
  return (
    <div className='mt-9 mx-auto'>
      <h1 className='text-white text-center font-extralight font-raleway'>CircleKu</h1>
      <p className={ errMsg === '' ? "offscreen" : "errMsg"}>{ errMsg }</p>
      <form onSubmit={ handleSubmit } className="mt-4">
        <div className="mb-2">
          <label
            htmlFor="username"
            className="block text-xs font-semibold"
          >
            Username
          </label>
          <input
            ref={usernameRef}
            autoComplete='off'
            id="username"
            type="text"
            className="block w-full px-4 py-1 mt-2 text-black bg-gray-dark rounded-md  focus:ring-purple-600 focus:outline-none focus:ring focus:ring-opacity-40"
          />
        </div>
        <div className="mb-2">
          <label
            htmlFor="email"
            className="block text-xs font-semibold mt-5"
          >
            Email
          </label>
          <input
            ref={emailRef}
            autoComplete='off'
            id="email"
            type="email"
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
            className="block w-full px-4 py-1 mt-2 text-black bg-gray-dark rounded-md  focus:ring-purple-600 focus:outline-none focus:ring focus:ring-opacity-40"
          />
        </div>
        <div className="mt-6 text-center">
          <button className="inline-block px-6 py-2 border-2 border-purple-600 text-purple-600 font-medium text-xs leading-tight uppercase rounded-full hover:bg-gray-800 focus:outline-none focus:ring-0 transition duration-150 ease-in-out">
            Signup
          </button>
        </div>
      </form>
      <p className="mt-5 text-xs font-light text-center text-gray-500">
        {" "}
        Already have an account?{" "}
        <a
          onClick = { props.handler }
          className="font-medium text-purple-600 hover:underline cursor-pointer"
        >
          login
        </a>
      </p>
    </div>
  )
}