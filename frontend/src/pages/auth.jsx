import dynamic from 'next/dynamic'

const Auth = dynamic(() => import('@/firebase/authHandler'), {
  ssr: false,
})

export default () => <Auth />