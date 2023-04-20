import React, { useState } from "react";
import { faXmark } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Montserrat } from "@next/font/google";

const montserrat = Montserrat({
  weight: "400",
  subsets: ["latin"],
});

const createRoom = async (name, picture, uid) => {
  const axios = require("axios");
  return axios.post(process.env.NEXT_PUBLIC_APP_HOST + "/api/private/create", {
    name: name,
    picture: picture,
    users: [{
      id: uid
    }],
    is_private: true
  }, {
    withCredentials: true,
  })
}

export default function GetUserModal({ setGetUsers, result, callback }) {
  const [query, setQuery] = useState("");

  const filtered = result.filter((user) => {
    return (
      user.Fullname.toLowerCase().includes(query.toLowerCase())
    );
  });

  const handleChange = (e) => {
    setQuery(e.target.value)
  };

  const handleCreateRoom = (name, picture, uid) => {
    createRoom(name, picture, uid)
      .then((response) => {
        setGetUsers(false);
        callback();
      })
      .catch((error) => {
        console.error(error)
      });
  }

  return (
    <div className="flex items-center justify-center fixed left-0 top-0 z-[1055] h-full w-full overflow-y-auto overflow-x-hidden backdrop-blur-sm bg-black/40 outline-none">
      <div className="h-96 w-1/2 max-w-3xl p-2 bg-[#222131] rounded-xl">
        <div className="flex gap-2 px-4 py-2 w-full rounded-md bg-[#2A293B]">
          <input
            type="text"
            name="fullname"
            style={montserrat}
            className="grow h-7 text-md bg-transparent text-gray-200 focus:outline-none"
            placeholder="Search contacts..."
            autoComplete="off"
            autoFocus
            onChange={handleChange}
          />
          <button onClick={() => setGetUsers(false)}>
            <FontAwesomeIcon
              className="text-2xl text-[#56537b]"
              icon={faXmark}>
            </FontAwesomeIcon>
          </button>
        </div>
        <ul className="">
          {filtered.map((data) => (
            <li
              key={data.ID}
              className="flex items-center gap-2 my-2 p-2 rounded-md hover:bg-[#2A293B] cursor-pointer"
              onClick={() => handleCreateRoom(data.Fullname, data.Picture, data.ID)}
            >
              <img
                className="rounded-full"
                src={data.Picture}
                width={32}
                height={32}
              />
              <span style={montserrat} className="text-md">{data.Fullname}</span>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
