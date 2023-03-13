import Head from 'next/head'
import Image from 'next/image'
import React from 'react'
import Auth from '@/components/auth'
import puplepattern from 'src/asset/texture-dark-background-purple-3840x2715-3086.jpg'

export default function Home() {
  return (
    <>
      <Head>
        <title>Sirkelin</title>
        <meta name="description" content="Sirkelin" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main>
        <div className='grid h-screen place-items-center' >
            <div className='bg-slate-900 rounded-xl shadow-2xl flex wh-l '>
                <Image src={puplepattern} className='wh-i rounded-xl' />
                <Auth />
          </div>
        </div>
      </main>
    </>
  )
}