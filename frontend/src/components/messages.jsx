import React, { useEffect, useState } from "react";
import { faPen, faSearch } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Montserrat } from "@next/font/google";
import Image from "next/image";
import GetUserModal from "./getUserModal";
import Loading from "@/components/loading";
import UserIcon from "src/asset/follower.png";

const montserrat = Montserrat({
  weight: "400",
  subsets: ["latin"],
});

const fetchRooms = async () => {
  const axios = require("axios");
  return axios.get(process.env.NEXT_PUBLIC_APP_HOST + "/api/private", {
    withCredentials: true,
  });
};

const mapRooms = (rooms) => {
  return rooms.map((room) => (
    <li
      className="h-20 px-4 flex items-center rounded-lg gap-2 cursor-pointer hover:bg-[#404040]/30"
      key={room.RoomId}
    >
      <Image src={UserIcon} className="w-10 h-10 " />
      {room.Username}
    </li>
  ));
};

export default function Messages() {
  const [rooms, setRooms] = useState(<Loading />);
  const [getUsers, setGetUsers] = useState(false);

  useEffect(() => {
    fetchRooms()
      .then((response) => {
        if (response.data.data.rooms === null) {
          setRooms(<span>Nothing to see here</span>);
        } else {
          setRooms(mapRooms(response.data.data.rooms));
        }
      })
      .catch((error) => {
        console.error(error);
      });
  }, []);

  return (
    <main>
      {getUsers ? <GetUserModal /> : null}
      <header>
        <h1 className="text-2xl text-center">Messages</h1>
      </header>
      <div className="w-96 px-4" style={montserrat.style}>
        <div className="">
          <section className="flex gap-3">
            <div className="flex items-center gap-3 py-2 px-4 rounded-full backdrop-blur-sm bg-[#323232]/30 hover:bg-[#404040]/30 focus-within:bg-[#404040]/30">
              <FontAwesomeIcon
                className="text-base text-gray-700"
                icon={faSearch}
              ></FontAwesomeIcon>
              <input
                style={montserrat.style}
                type="text"
                className="grow h-7 text-sm bg-transparent text-gray-200 focus:outline-none"
                placeholder="Search..."
              />
            </div>
            <button
              className="w-12 p-2 rounded-full text-indigo-700 border border-indigo-700 hover:bg-indigo-700 hover:text-gray-200"
              onClick={() => setGetUsers(true)}
            >
              <FontAwesomeIcon className="" icon={faPen}></FontAwesomeIcon>
            </button>
          </section>
          <ul className="my-6">{rooms}</ul>
        </div>
      </div>
    </main>
  );
}
