import Head from "next/head";
import { useRouter } from "next/router";
import Link from "next/link";
import Image from "next/image";
import { Montserrat, Yantramanav } from "@next/font/google";
import React, { useState, useEffect } from "react";
import { UserAuth } from "@/context/AuthContext";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faComment,
  faHouse,
  faUserGroup,
  faInfinity,
  faArrowRightFromBracket,
} from "@fortawesome/free-solid-svg-icons";
import Messages from "./messages";
import Circle from "./circle";
import UserIcon from "@/asset/follower.png";

const yantramanav = Yantramanav({
  weight: "400",
  subsets: ["latin"],
});

const yBold = Yantramanav({
  weight: "500",
  subsets: ["latin"],
});

const montserrat = Montserrat({
  weight: "400",
  subsets: ["latin"],
});

export default function Dashboard() {
  const router = useRouter();
  const { page } = router.query;
  const [view, setView] = useState(null);
  const [selector, setSelector] = useState(1);
  const { firebaseSignOut, user } = UserAuth();

  useEffect(() => {
    switch (page) {
      case "messages":
        setView(<Messages />);
        setSelector(1);
        break;
      case "circle":
        setView(<Circle />);
        setSelector(2);
        break;
      default:
        <ErrorPage statusCode={404} />; //TODO: ganti ke error page
        break;
    }
  }, [page]);
  const title = page.charAt(0).toUpperCase() + page.slice(1);

  return (
    <>
      <Head>
        <title>{title}</title>
        <meta name="description" content="Generated by create next app" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
      </Head>
      <main>
        <div className="h-screen">
          <div className="wh-chat flex">
            <div className="h-screen flex flex-col navbar-ct backdrop-blur-sm bg-[#323232]/30">
              <header>
                <h1
                  className="text-3xl text-center leading-normal translate-y-0.5 grow"
                  style={montserrat.style}
                >
                  Sirkelin.
                </h1>
              </header>
              <div className="grow">
                <nav className="px-4 grid text-center gap-8 text-xl">
                  <Link href="/messages" legacyBehavior>
                    <button className={"flex group hover:bg-gray-700/40 rounded-xl py-3 px-4 " + (selector === 1 ? "bg-gray-700/40" : "")}>
                      <FontAwesomeIcon
                        className="self-center group-hover:stroke-white"
                        icon={faComment}
                      />
                      <span
                        style={yantramanav.style}
                        className="ml-4 font-yantramanav"
                      >
                        messages
                      </span>
                    </button>
                  </Link>
                  <Link href="/circle" legacyBehavior>
                    <button className={"flex group hover:bg-gray-700/40 rounded-xl py-3 px-4 " + (selector === 2 ? "bg-gray-700/40" : "")}>
                      <FontAwesomeIcon
                        className="self-center group-hover:stroke-white"
                        icon={faInfinity}
                      />
                      <span
                        style={yantramanav.style}
                        className=" ml-3 font-yantramanav"
                      >
                        circle
                      </span>
                    </button>
                  </Link>
                  <button className={"flex group hover:bg-gray-700/40 rounded-xl py-3 px-4 " + (selector === 3 ? "bg-gray-700/40" : "")}>
                    <FontAwesomeIcon
                      className="self-center group-hover:stroke-white"
                      icon={faUserGroup}
                    />
                    <span
                      style={yantramanav.style}
                      className=" ml-3 font-yantramanav"
                    >
                      friend list
                    </span>
                  </button>
                  <button className={"flex group hover:bg-gray-700/40 rounded-xl py-3 px-4 " + (selector === 4 ? "bg-gray-700/40" : "")}>
                    <FontAwesomeIcon
                      className="self-center group-hover:stroke-white"
                      icon={faHouse}
                    />
                    <span
                      style={yantramanav.style}
                      className=" ml-3 font-yantramanav"
                    >
                      explore
                    </span>
                  </button>
                </nav>
              </div>
              <div className="m-4 p-2 rounded-lg">
                <div className="flex items-center gap-4">
                  <Image
                    alt="user photo"
                    className="rounded-full"
                    src={(user?.photoURL || UserIcon)}
                    width={36}
                    height={36}
                  />
                  <span className="w-full">
                    <b className="text-lg" style={yBold.style}>
                      {user?.displayName}
                    </b>
                  </span>
                  <FontAwesomeIcon
                    className="text-lg text-red-800 cursor-pointer"
                    icon={faArrowRightFromBracket}
                    onClick={firebaseSignOut}
                  />
                </div>
              </div>
            </div>
            {view}
          </div>
        </div>
      </main>
    </>
  );
}
