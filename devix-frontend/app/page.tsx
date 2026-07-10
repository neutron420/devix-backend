"use client";

import { HeroSection03 } from "@/components/hero-03";
import { Component as HeroSection } from "@/components/hero";
import FlowArt, { FlowSection } from "@/components/story-scroll";
import RuixenBentoCards from "@/components/ruixen-bento-cards";
import GlobeFeatureSection from "@/components/globe-feature-section";
import FAQs from "@/components/text-reveal-faqs";
import { AnimatedTestimonials } from "@/components/animated-testimonials";
import Footer from "@/components/footer";

const mockTestimonials = [
  {
    id: 1,
    name: "Alex Rivera",
    role: "Senior Full Stack Developer",
    company: "TechFlow",
    content: "Devix has completely transformed how our engineering team shares knowledge. The real-time discussion features are a game-changer.",
    rating: 5,
    avatar: "https://i.pravatar.cc/150?u=a042581f4e29026024d"
  },
  {
    id: 2,
    name: "Sarah Chen",
    role: "Engineering Manager",
    company: "GlobalNet",
    content: "The markdown support combined with direct messaging makes this the best technical collaboration platform I've used in my 10-year career.",
    rating: 5,
    avatar: "https://i.pravatar.cc/150?u=a042581f4e29026704d"
  },
  {
    id: 3,
    name: "David Kim",
    role: "Lead DevOps",
    company: "CloudScale",
    content: "Finally, a platform that understands what developers need. The build logs and organization workspaces are brilliantly implemented.",
    rating: 5,
    avatar: "https://i.pravatar.cc/150?u=a04258114e29026702d"
  }
];

const mockCompanies = ["TechFlow", "GlobalNet", "CloudScale", "DevCorp", "InnovateTech"];

export default function Home() {
  return (
    <div className="w-full min-h-screen flex flex-col bg-background">
      <HeroSection03 />
      <HeroSection />
      <FlowArt aria-label="Devix Platform Presentation">
        <FlowSection aria-label="What is Devix?" style={{ backgroundColor: '#fd5200', color: '#fff' }}>
          <p className="text-xs font-bold uppercase tracking-[0.2em]">01 — What is Devix?</p>
          <hr className="my-[2vw] border-none border-t border-black opacity-100" />
          <div>
            <h1
              className="text-[clamp(3.5rem,12vw,14rem)] font-bold leading-[0.85] uppercase tracking-tight"
            >
              Share.
              <br />
              Discuss.
              <br />
              Build.
            </h1>
          </div>
          <hr className="my-[2vw] border-none border-t border-black opacity-100" />
          <p className="mt-auto max-w-[50ch] text-[clamp(1rem,2.5vw,2rem)] font-normal leading-relaxed">
            Devix is a modern, interactive, and real-time social platform built for developers to share technical knowledge, collaborate, and network. It bridges the gap between technical blogging, real-time collaboration, and community discussion.
          </p>
        </FlowSection>

        <FlowSection aria-label="Feeds & Posts" style={{ backgroundColor: '#000', color: '#fff' }}>
          <p className="text-xs font-bold uppercase tracking-[0.2em]">02 — Feeds & Posts</p>
          <hr className="my-[2vw] border-none border-t border-white/60" />
          <div>
            <h2
              className="text-[clamp(3.5rem,12vw,14rem)] font-bold leading-[0.85] uppercase tracking-tight"
            >
              Post.
              <br />
              Vote.
              <br />
              Discuss.
            </h2>
          </div>
          <hr className="my-[2vw] border-none border-t border-white/60" />
          <p className="max-w-[50ch] text-[clamp(1rem,2.5vw,2rem)] font-normal leading-relaxed">
            Publish articles using a full markdown editor with image upload support. Browse through the global feed or view a personalized stream from developers you follow.
          </p>
          <hr className="my-[2vw] border-none border-t border-white/60" />
          <div className="flex flex-wrap gap-[3vw]">
            <div className="min-w-[180px] flex-1">
              <p className="mb-2 text-sm font-bold uppercase tracking-wider">Post Types</p>
              <p className="text-[clamp(0.85rem,1.3vw,1.05rem)] leading-relaxed opacity-75">
                Publish articles categorized as Questions, Concepts, or Build Logs with full markdown support.
              </p>
            </div>
            <div className="min-w-[180px] flex-1">
              <p className="mb-2 text-sm font-bold uppercase tracking-wider">Feeds</p>
              <p className="text-[clamp(0.85rem,1.3vw,1.05rem)] leading-relaxed opacity-75">
                Browse through the global feed (Latest/Trending) or view a personalized stream of followed developers.
              </p>
            </div>
            <div className="min-w-[180px] flex-1">
              <p className="mb-2 text-sm font-bold uppercase tracking-wider">Interactivity</p>
              <p className="text-[clamp(0.85rem,1.3vw,1.05rem)] leading-relaxed opacity-75">
                Create interactive polls, cast votes, upvote/downvote posts, and participate in threaded comments.
              </p>
            </div>
          </div>
        </FlowSection>

        <FlowSection aria-label="Real-time Ecosystem" style={{ backgroundColor: '#F5F0E8', color: '#000' }}>
          <p className="text-xs font-bold uppercase tracking-[0.2em]">03 — Real-Time features</p>
          <hr className="my-[2vw] border-none border-t border-black/60" />
          <div>
            <h2
              className="text-[clamp(3.5rem,12vw,14rem)] font-bold leading-[0.85] uppercase tracking-tight"
            >
              Chat.
              <br />
              Notify.
              <br />
              Engage.
            </h2>
          </div>
          <hr className="my-[2vw] border-none border-t border-black/60" />
          <p className="max-w-[50ch] text-[clamp(1rem,2.5vw,2rem)] font-normal leading-relaxed">
            Real-time direct messaging and instant notification updates keep the ecosystem alive and interactive.
          </p>
          <hr className="my-[2vw] border-none border-t border-black/60" />
          <div className="flex flex-wrap gap-[3vw]">
            <div className="min-w-[180px] flex-1">
              <p className="mb-2 text-sm font-bold uppercase tracking-wider">Direct Messaging</p>
              <p className="text-[clamp(0.85rem,1.3vw,1.05rem)] leading-relaxed opacity-75">
                Instant private chats with online/offline presence indicators, live typing indicators, and read receipts.
              </p>
            </div>
            <div className="min-w-[180px] flex-1">
              <p className="mb-2 text-sm font-bold uppercase tracking-wider">Notifications</p>
              <p className="text-[clamp(0.85rem,1.3vw,1.05rem)] leading-relaxed opacity-75">
                Stay updated with instant alerts for new comments, likes, follower updates, and incoming chat messages.
              </p>
            </div>
            <div className="min-w-[180px] flex-1">
              <p className="mb-2 text-sm font-bold uppercase tracking-wider">Bookmarks & Logs</p>
              <p className="text-[clamp(0.85rem,1.3vw,1.05rem)] leading-relaxed opacity-75">
                Save posts to read later and view your personal activity history to track your contributions.
              </p>
            </div>
          </div>
        </FlowSection>

        <FlowSection aria-label="Analytics & Organizations" style={{ backgroundColor: '#1A3DE8', color: '#fff' }}>
          <p className="text-xs font-bold uppercase tracking-[0.2em]">04 — Teams & Stats</p>
          <hr className="my-[2vw] border-none border-t border-white/50" />
          <div>
            <h2
              className="text-[clamp(3.5rem,12vw,14rem)] font-bold leading-[0.85] uppercase tracking-tight"
            >
              Analyze.
              <br />
              Group.
              <br />
              Grow.
            </h2>
          </div>
          <hr className="my-[2vw] border-none border-t border-white/50" />
          <p className="max-w-[50ch] text-[clamp(1rem,2.5vw,2rem)] font-normal leading-relaxed">
            Gain insights on your technical writing and form organizations to collaborate with developers in team spaces.
          </p>
          <hr className="my-[2vw] border-none border-t border-white/50" />
          <div className="flex flex-wrap gap-[3vw]">
            <div className="min-w-[180px] flex-1">
              <p className="mb-2 text-sm font-bold uppercase tracking-wider">Author Analytics</p>
              <p className="text-[clamp(0.85rem,1.3vw,1.05rem)] leading-relaxed opacity-75">
                Track post performance using dashboard analytics showing view counts, referrers, and demographics.
              </p>
            </div>
            <div className="min-w-[180px] flex-1">
              <p className="mb-2 text-sm font-bold uppercase tracking-wider">Organizations</p>
              <p className="text-[clamp(0.85rem,1.3vw,1.05rem)] leading-relaxed opacity-75">
                Form groups and organizations to build dedicated spaces for team collaboration and shared goals.
              </p>
            </div>
            <div className="min-w-[180px] flex-1">
              <p className="mb-2 text-sm font-bold uppercase tracking-wider">Developer Home</p>
              <p className="text-[clamp(0.85rem,1.3vw,1.05rem)] leading-relaxed opacity-75">
                A specialized engineering knowledge platform designed to host developer-centric activities.
              </p>
            </div>
          </div>
        </FlowSection>

        <FlowSection aria-label="Join us" style={{ backgroundColor: '#000', color: '#fff' }}>
          <p className="text-xs font-bold uppercase tracking-[0.2em]">05 — Join us</p>
          <hr className="my-[2vw] border-none border-t border-black/60" />
          <div>
            <h2
              className="text-[clamp(3.5rem,12vw,14rem)] font-bold leading-[0.85] uppercase tracking-tight"
            >
              Ready
              <br />
              To
              <br />
              Build?
            </h2>
          </div>
          <hr className="my-[2vw] border-none border-t border-black/60" />
          <p className="mt-auto max-w-[50ch] text-[clamp(1rem,2.5vw,2rem)] font-normal leading-relaxed">
            Take control of your engineering journey. Join Devix now and let&apos;s shape the future of technical knowledge sharing together.
          </p>
        </FlowSection>
      </FlowArt>
      <RuixenBentoCards />
      <GlobeFeatureSection />
      <FAQs />
      <AnimatedTestimonials testimonials={mockTestimonials} trustedCompanies={mockCompanies} />
      <Footer />
    </div>
  );
}

