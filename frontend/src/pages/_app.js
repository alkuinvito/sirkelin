import '@/styles/globals.css'
import { config } from '@fortawesome/fontawesome-svg-core'
import '@fortawesome/fontawesome-svg-core/styles.css'
import React, { useState } from 'react'
import {LoginContext} from '@/components/Context.js'
config.autoAddCss = false

export default function App({ Component, pageProps }) {
  const [info, setInfo] = useState({username: '', picture: ''})
  return(
    <LoginContext.Provider value={[info, setInfo]}>
      <Component {...pageProps}/>
    </LoginContext.Provider>
  )
}