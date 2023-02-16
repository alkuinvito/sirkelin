import React, { useEffect } from 'react'
import firebase from '@/firebase/clientApp'

const uiConfig = {
  callbacks: {
    signInSuccessWithAuthResult: () => true,
    uiShown: function() {
      // TODO: hide loader here
    }
  },
  signInFlow: 'popup',
  signInSuccessUrl: 'http://localhost:3000/',
  signInOptions: [
    firebase.auth.GoogleAuthProvider.PROVIDER_ID
  ]
}

export default function authHandler() {
  useEffect(() => {
    const ui = firebaseui.auth.AuthUI.getInstance() || new firebaseui.auth.AuthUI(firebase.auth())
    ui.start('#firebaseui-container', uiConfig)
  })

  return (
    <section>
      <h1 className='text-2xl'>Please sign-in:</h1>
      <div id='firebaseui-container'></div>
    </section>
  )
}