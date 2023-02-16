import React from 'react'
import { UserAuth } from '@/context/AuthContext'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faGoogle } from '@fortawesome/free-brands-svg-icons'

function Auth() {
  const { firebaseSignIn, user } = UserAuth()
  const signInHandler = async () => {
    try {
      await firebaseSignIn()
    } catch(error) {
      console.log(error)
    }
  }

  return (
    <div>
      <h1 className='text-gray-400'>Please sign in</h1>
      <button className='p-2 rounded-md bg-[#4285F4] hover:bg-[#4285F4]/80 transition' onClick={signInHandler}><FontAwesomeIcon className='mr-2' icon={faGoogle} />Sign in with Google</button>
    </div>
  )
}

export default Auth