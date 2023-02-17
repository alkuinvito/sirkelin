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

const AuthContext = createContext()

const createSession = async (url, idToken, csrfToken) => {
  const {data} = await Axios.post(url, {idToken, csrfToken})
  return data
}

export const AuthContextProvider = ({ children }) => {
  const [ user, setUser ] = useState({})

  const firebaseSignIn = () => {
    const provider = new GoogleAuthProvider()
    setPersistence(auth, inMemoryPersistence)
    signInWithPopup(auth, provider)
      .then(userCredential => {
        return userCredential.user.getIdToken().then(idToken => {
          const csrfToken = getCookie('csrfToken')
          return createSession('/api/user/sessionLogin/', idToken, csrfToken)
        })
      })
      .then(result => {
        setUser(result.user)
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