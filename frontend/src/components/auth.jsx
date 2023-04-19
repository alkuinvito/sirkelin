import React from "react";
import { UserAuth } from "@/context/AuthContext";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faGoogle, faGithub } from "@fortawesome/free-brands-svg-icons";
import { GoogleAuthProvider, GithubAuthProvider } from "firebase/auth";

function Auth() {
  const googleProvider = new GoogleAuthProvider();
  const githubProvider = new GithubAuthProvider();
  const { firebaseSignIn } = UserAuth();
  const signInHandler = async (provider) => {
    try {
      await firebaseSignIn(provider);
    } catch (error) {
      console.log(error);
    }
  };

  return (
    <div className="p-8">
      <h1 className="text-gray-200 text-2xl font-normal mb-8">Sign in to <b className="bg-clip-text text-transparent bg-gradient-to-tr from-indigo-700 to-pink-400">Sirkelin</b> using social account</h1>
      <button
        className="w-full p-2 mb-4 rounded-md bg-[#4285F4] hover:bg-[#4285F4]/80"
        onClick={() => signInHandler(googleProvider)}
      >
        <FontAwesomeIcon className="mr-2" icon={faGoogle} />
        Sign in with Google
      </button>
      <button
        className="w-full p-2 rounded-md bg-[#262627] hover:bg-[#262627]/80"
        onClick={() => signInHandler(githubProvider)}
      >
        <FontAwesomeIcon className="mr-2" icon={faGithub} />
        Sign in with Github
      </button>
    </div>
  );
}

export default Auth;
