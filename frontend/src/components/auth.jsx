import React from 'react'
import { UserAuth } from '@/context/AuthContext'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faGoogle, faGithub } from '@fortawesome/free-brands-svg-icons'
import { GoogleAuthProvider, GithubAuthProvider } from 'firebase/auth'

function Auth() {
  const googleProvider = new GoogleAuthProvider()
  const githubProvider = new GithubAuthProvider()
  const { firebaseSignIn } = UserAuth()
  const signInHandler = async (provider) => {
    try {
      await firebaseSignIn(provider)
    } catch(error) {
      console.log(error)
    }
  }

  return (
    <div>
      <h1 className='text-gray-400'>Please sign in</h1>
      <button className='p-2 rounded-md bg-[#4285F4] hover:bg-[#4285F4]/80 transition' onClick={() => signInHandler(googleProvider)}><FontAwesomeIcon className='mr-2' icon={faGoogle} />Sign in with Google</button>
      <button className='p-2 rounded-md bg-[#262627] hover:bg-[#262627]/80 transition' onClick={() => signInHandler(githubProvider)}><FontAwesomeIcon className='mr-2' icon={faGithub} />Sign in with Github</button>      
    </div>
  )
}

export default Auth