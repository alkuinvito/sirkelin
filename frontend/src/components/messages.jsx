import React, { useEffect, useState } from "react";
import { faPen, faSearch } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Montserrat } from "@next/font/google";
import Image from "next/image";
import GetUserModal from "./getUserModal";
import Loading from "@/components/loading";
import Room from "./room";

const montserrat = Montserrat({
  weight: "400",
  subsets: ["latin"],
});

const fetchUsers = async () => {
  const axios = require("axios");
  return axios.get(process.env.NEXT_PUBLIC_APP_HOST + "/api/user/list");
};

const fetchRooms = async () => {
  const axios = require("axios");
  return axios.get(process.env.NEXT_PUBLIC_APP_HOST + "/api/private");
};

const mapRooms = (rooms, setter) => {
  return rooms.map((room) => (
    <li
      className="h-16 px-4 flex gap-4 items-center rounded-lg cursor-pointer hover:bg-gray-700/30"
      key={room.RoomId}
      onClick={() => { setter(room.RoomId, room.Fullname, room.Picture) }}
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
  const [room, setRoom] = useState({});
  const [getUsers, setGetUsers] = useState(false);
  const [result, setResult] = useState([]);

  const handleChangeRoom = (RoomId, Name, Picture) => {
    setRoom({
      RoomId: RoomId,
      Name: Name,
      Picture: Picture
    });
  };

  const updateRooms = () => {
    fetchRooms()
      .then((response) => {
        if (response.data.data.rooms === null) {
          setRooms(<span>Nothing to see here</span>);
        } else {
          setRooms(mapRooms(response.data.data.rooms, handleChangeRoom));
        }
      })
      .catch((error) => {
        console.error(error);
      });
  }

  useEffect(() => {
    updateRooms();
    fetchUsers()
      .then((response) => {
        if (response.data.data.users === null) {
          setResult([]);
        } else {
          setResult(response.data.data.users);
        }
      })
      .catch(console.error);

    const refreshRoom = setInterval(() => {
      updateRooms();
    }, 12000);

    return () => clearInterval(refreshRoom);
  }, []);

  return (
    <main className="flex grow">
      {getUsers ? <GetUserModal setGetUsers={setGetUsers} result={result} callback={updateRooms} /> : null}
      <section>
        <header>
          <h1 className="text-2xl text-center">Messages</h1>
        </header>
        <div className="w-80 px-4" style={montserrat.style}>
          <div>
            <div className="flex gap-3">
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
            </div>
            <ul className="my-6">{rooms}</ul>
          </div>
        </div>
      </section>
      <section className="grow border-l border-gray-700/40">
        {room?.RoomId === undefined ? null : <Room room={room} />}
      </section>
    </main>
  );
}
