import React, { useEffect, useState } from "react";
import { faPen, faSearch } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Montserrat } from "@next/font/google";
import Image from "next/image";
import GetUserModal from "./getUserModal";
import Loading from "@/components/loading";

const montserrat = Montserrat({
  weight: "400",
  subsets: ["latin"],
});

const fetchUsers = async () => {
  const axios = require("axios");
  return axios.get(process.env.NEXT_PUBLIC_APP_HOST + "/api/user/list", {
    withCredentials: true,
  });
};

const fetchRooms = async () => {
  const axios = require("axios");
  return axios.get(process.env.NEXT_PUBLIC_APP_HOST + "/api/private", {
    withCredentials: true,
  });
};

const mapRooms = (rooms) => {
  return rooms.map((room) => (
    <li
      className="h-16 px-4 flex gap-4 items-center rounded-lg cursor-pointer hover:bg-gray-700/30"
      key={room.RoomId}
    >
      <Image
        alt="contact photo"
        className="rounded-full"
        src={room.Picture}
        width={32}
        height={32}
      />
      <span style={montserrat}>{room.Fullname}</span>
    </li>
  ));
};

export default function Messages() {
  const [rooms, setRooms] = useState(<Loading />);
  const [getUsers, setGetUsers] = useState(false);
  const [result, setResult] = useState([]);

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
    fetchUsers()
      .then((response) => {
        if (response.data.data.users === null) {
          setResult([]);
        } else {
          setResult(response.data.data.users);
        }
      })
      .catch(console.error);
  }, []);

  return (
    <main>
      {getUsers ? <GetUserModal setGetUsers={setGetUsers} result={result} /> : null}
      <header>
        <h1 className="text-2xl text-center">Messages</h1>
      </header>
      <div className="w-80 px-4" style={montserrat.style}>
        <div>
          <section className="flex gap-3">
            <div className="flex items-center gap-3 grow py-2 px-4 rounded-full backdrop-blur-sm bg-gray-700/40 hover:bg-gray-700/30 focus-within:bg-gray-700/30">
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
              className="text-gray-700 hover:text-indigo-700"
              onClick={() => setGetUsers(!getUsers)}
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
