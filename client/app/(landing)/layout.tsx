import React from "react";
import Navbar from "./homepage/Navbar";
import Footer from "./homepage/Footer";

interface Children {
  children: React.ReactNode;
}

const HomepageLayout = ({ children }: Children) => {
  return (
    <div className="bg-[url('/assets/bg.png')] bg-cover bg-center bg-no-repeat">
      <Navbar />
      <main>{children}</main>
      <Footer />
    </div>
  );
};

export default HomepageLayout;
