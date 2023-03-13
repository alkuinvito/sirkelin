import { useContext, createContext, useState } from 'react'
import {
  setPersistence,
  inMemoryPersistence,
  GoogleAuthProvider,
  signInWithPopup,
  signOut
} from 'firebase/auth'
import { auth } from '@/firebase/clientApp'
import { Axios } from 'axios'
import { useRouter } from 'next/router'

const AuthContext = createContext()

const createSession = async (url, idToken, csrfToken) => {
  const {data} = await Axios.post(url, {idToken, csrfToken})
  return data
}

export const AuthContextProvider = ({ children }) => {
  const [ user, setUser ] = useState({})
  const router = useRouter()

  const firebaseSignIn = () => {
    const provider = new GoogleAuthProvider()
    setPersistence(auth, inMemoryPersistence)
    signInWithPopup(auth, provider)
      .then(userCredential => {
        return userCredential.user.getIdToken().then(idToken => {
          setUser({
            displayName: userCredential.user.displayName,
            photoURL: userCredential.user.photoURL
          })
          console.log(idToken)
          router.push('/messages')
          // const csrfToken = getCookie('csrfToken')
          // return createSession('/api/user/sessionLogin/', idToken, csrfToken)
        })
      })
      .catch(error => {
        const errorCode = error.code
        const errorMessage = error.message
        console.error(errorCode, errorMessage)
      })
  }

  const firebaseSignOut = () => {
    signOut(auth)
  }

  return (
    <AuthContext.Provider value={{ firebaseSignIn, firebaseSignOut, user }}>
      {children}
    </AuthContext.Provider>
  )
}

export const UserAuth = () => {
  return useContext(AuthContext)
}