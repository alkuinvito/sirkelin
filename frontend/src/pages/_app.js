import '@/styles/globals.css'
import { config } from '@fortawesome/fontawesome-svg-core'
import '@fortawesome/fontawesome-svg-core/styles.css'
import { AuthContextProvider } from '@/context/AuthContext'
config.autoAddCss = false

export default function App({ Component, pageProps }) {
  return(
    <AuthContextProvider>
      <Component {...pageProps}/>
    </AuthContextProvider>
  )
}