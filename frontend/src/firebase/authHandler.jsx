import React, { useEffect } from 'react'
import firebase from '@/firebase/clientApp'

const firebaseui = require('firebaseui')

export default function authHandler() {
  useEffect(() => {
    const ui = firebaseui.auth.AuthUI.getInstance() || new firebaseui.auth.AuthUI(firebase.auth())
    ui.start('#firebaseui-container', {
      signInOptions: [
        firebase.auth.GoogleAuthProvider.PROVIDER_ID
      ],      
    }, []);
  })

  return (
    <div id='firebaseui-container'>
    </div>
  )
}