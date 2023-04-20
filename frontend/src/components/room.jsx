import React, { useState, useEffect } from "react";
import Image from "next/image";
import { useLocalStorage } from "@/context/useLocalStorage";
import { Montserrat } from "@next/font/google";
import { faPaperPlane } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

const montserrat = Montserrat({
  weight: "400",
  subsets: ["latin"],
});

const fetchRoom = async (RoomId) => {
  const axios = require("axios");
  return axios.get(process.env.NEXT_PUBLIC_APP_HOST + "/api/room/" + RoomId)
}

const sendMessage = async (RoomId, body) => {
  const axios = require("axios");
  return axios.post(process.env.NEXT_PUBLIC_APP_HOST + "/api/room/" + RoomId, {
    body: body
  })
}

const mapMessages = (messages, uid) => {
  return (
    messages.map((message) => (
      <li
        style={montserrat.style}
        key={message.ID}
        className={"py-2 px-3 mb-2 w-fit max-w-sm box-border rounded-xl" + (message.UserID === uid ? " justify-self-end bg-indigo-800/40" : " bg-gray-500/40")}
      >
        <span>{message.Body}</span>
      </li>
    ))
  )
}

export default function Room(props) {
  const [messages, setMessages] = useState([]);
  const [message, setMessage] = useState("");
  const [user, _] = useLocalStorage("user");

  const handleSubmit = (e) => {
    e.preventDefault();
    console.log(message)
    if (message.length !== 0) {
      setMessage("");
      sendMessage(props.room?.RoomId, message)
        .then((response) => {
        })
        .catch((error) => {
          console.error(error);
        });
    }
  }

  useEffect(() => {
    if (props.room?.RoomId !== undefined) {
      fetchRoom(props.room?.RoomId)
        .then((response) => {
          const responseMsg = response.data.data.messages;
          if (messages !== responseMsg) {
            console.log("updating messages");
            setMessages(mapMessages(responseMsg, user.ID));
          }
        })
        .catch((error) => {
          console.error(error);
        })
      const refreshMsg = setInterval(() => {
        fetchRoom(props.room?.RoomId)
          .then((response) => {
            const responseMsg = response.data.data.messages;
            if (messages !== responseMsg) {
              console.log("updating messages");
              setMessages(mapMessages(responseMsg, user.ID));
            }
          })
          .catch((error) => {
            console.error(error);
          })
      }, 3000);

      return () => clearInterval(refreshMsg);
    }
  }, [props.room]);

  return (
    <section className="h-screen grow border-l border-gray-700/40" style={montserrat.style}>
      <header className="border-b border-gray-700/40">
        <Image
          src={props.room?.Picture}
          width={40}
          height={40}
        >
        </Image>
        <h1>{props.room?.Name}</h1>
      </header>
      <div className="ct-room flex flex-col h-full">
        <div className="p-4 grow">
          <ol className="grid">
            {messages}
          </ol>
        </div>
        <form className="flex items-center w-full p-3 border-t border-gray-700/40" onSubmit={(e) => { handleSubmit(e) }}>
          <input
            style={montserrat.style}
            className="w-full py-2 px-4 bg-gray-700/40 text-md rounded-md hover:bg-gray-700/30 focus-within:bg-gray-700/30 focus:outline-none focus:ring-1 focus:ring-indigo-800"
            type="text"
            name="body"
            autoComplete="off"
            value={message}
            onChange={(e) => { setMessage(e.target.value) }}
          />
          <button className="pl-4 pr-3 text-xl text-gray-700 hover:text-indigo-700">
            <FontAwesomeIcon icon={faPaperPlane} />
          </button>
        </form>
      </div>
    </section>
  )
}