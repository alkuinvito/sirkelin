"use client";

import React from "react";
import { Button, Space } from "antd";
import { GoogleAuthProvider, GithubAuthProvider } from "firebase/auth";
import { useAuth } from "@/contexts/AuthContext";

export default function Page() {
  const google = new GoogleAuthProvider();
  const github = new GithubAuthProvider();
  const { firebaseSignIn } = useAuth();
  const handleSignIn = async (provider) => {
    try {
      await firebaseSignIn(provider);
    } catch (error) {
      console.log(error);
    }
  };

  return (
    <div style={{ padding: "0 24px" }}>
      <Space size="middle" style={{ display: "flex" }}>
        <Button type="primary" onClick={() => handleSignIn(google)}>
          Sign in with Google
        </Button>
        <Button type="primary" onClick={() => handleSignIn(github)}>
          Sign in with Github
        </Button>
      </Space>
    </div>
  );
}
