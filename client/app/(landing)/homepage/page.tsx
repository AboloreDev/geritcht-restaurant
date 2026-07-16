import React from "react";
import HomepageLayout from "../layout";
import Hero from "./Hero";
import About from "./About";
import SpecialMenu from "./SpecialMenu";
import Chef from "./Chef";
import VideoIntro from "./Video";
import Awards from "./Awards";
import Gallery from "./Gallery";
import Contact from "./Contact";
import Newsletter from "./Newsletter";

const Home = () => {
  return (
    <>
      <div>
        <Hero />
        <About />
        <SpecialMenu />
        <Chef />
        <VideoIntro />
        <Awards />
        <Gallery />
        <Contact />
        <Newsletter />
      </div>
    </>
  );
};

export default Home;
