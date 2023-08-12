"use client";

import RoomList from "@/components/roomList";
import { useAuth } from "@/contexts/AuthContext";

export default function Page() {
  const { user, firebaseSignOut } = useAuth();

  return (
    <div>
      <h1>{user?.displayName}</h1>
      <RoomList />
    </div>
  );
}
