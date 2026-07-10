"use client";

import { motion } from "framer-motion";
import {
  FiZap,
  FiUsers,
  FiArrowUp,
  FiMessageCircle,
  FiBookmark,
  FiShare2,
  FiHome,
  FiCompass,
  FiPlusCircle,
  FiBell,
  FiUser,
  FiCode,
  FiSend,
  FiCheck,
  FiSearch,
} from "react-icons/fi";
import { SiNextdotjs, SiTailwindcss } from "react-icons/si";
import { Avatar, AvatarImage } from "@/components/ui/avatar";

export function AuthVisual() {
  return (
    <div className="absolute inset-0 flex flex-col items-center justify-center bg-gradient-to-br from-blue-600 via-sky-500 to-cyan-400 text-white select-none overflow-hidden h-full w-full">

      {/* ── Decorative Background ── */}

      {/* Diagonal Slash Lines */}
      <div className="absolute inset-0 pointer-events-none opacity-[0.07]">
        <div className="absolute top-0 left-0 w-full h-full" style={{
          backgroundImage: `repeating-linear-gradient(
            -45deg,
            transparent,
            transparent 40px,
            white 40px,
            white 42px
          )`
        }} />
      </div>

      {/* Grid Dots */}
      <div className="absolute inset-0 bg-[radial-gradient(circle,_rgba(255,255,255,0.12)_1px,_transparent_1px)] bg-[size:24px_24px] pointer-events-none" />

      {/* Big Ambient Glows */}
      <div className="absolute -top-20 -right-20 w-[500px] h-[500px] rounded-full bg-white/10 blur-[100px] pointer-events-none" />
      <div className="absolute -bottom-20 -left-20 w-[400px] h-[400px] rounded-full bg-sky-300/20 blur-[80px] pointer-events-none" />
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[300px] h-[300px] rounded-full bg-blue-400/15 blur-[60px] pointer-events-none" />

      {/* Floating Geometric Shapes */}
      <motion.div
        animate={{ rotate: [0, 360] }}
        transition={{ repeat: Infinity, duration: 30, ease: "linear" }}
        className="absolute top-[10%] right-[8%] w-16 h-16 border-2 border-white/10 rounded-xl pointer-events-none"
      />
      <motion.div
        animate={{ rotate: [0, -360] }}
        transition={{ repeat: Infinity, duration: 25, ease: "linear" }}
        className="absolute bottom-[15%] left-[5%] w-12 h-12 border-2 border-white/10 rounded-full pointer-events-none"
      />
      <motion.div
        animate={{ rotate: [45, 405] }}
        transition={{ repeat: Infinity, duration: 20, ease: "linear" }}
        className="absolute top-[60%] right-[5%] w-8 h-8 border-2 border-white/8 rounded-md pointer-events-none"
      />

      {/* ── Phone + Widgets Container ── */}
      <div className="relative flex items-center justify-center w-full flex-1 max-w-lg scale-[0.82] xl:scale-[0.88] 2xl:scale-[0.92]">

        {/* Phone Side Buttons */}
        <div className="absolute top-[calc(50%-155px)] -translate-x-[141px] w-[3px] h-7 bg-slate-700 rounded-l-sm z-10" />
        <div className="absolute top-[calc(50%-110px)] -translate-x-[141px] w-[3px] h-10 bg-slate-700 rounded-l-sm z-10" />
        <div className="absolute top-[calc(50%-50px)] -translate-x-[141px] w-[3px] h-10 bg-slate-700 rounded-l-sm z-10" />
        <div className="absolute top-[calc(50%-95px)] translate-x-[140px] w-[3px] h-14 bg-slate-700 rounded-r-sm z-10" />

        {/* Phone Frame */}
        <motion.div
          initial={{ y: 30, opacity: 0 }}
          animate={{ y: [0, -10, 0], opacity: 1 }}
          transition={{
            y: { repeat: Infinity, duration: 6, ease: "easeInOut" },
            opacity: { duration: 0.8 },
          }}
          className="relative w-[280px] h-[570px] rounded-[46px] border-[6px] border-slate-900 bg-white shadow-[0_30px_80px_-20px_rgba(0,0,0,0.4)] overflow-hidden flex flex-col z-20"
        >
          {/* ── iOS Status Bar ── */}
          <div className="absolute top-0 inset-x-0 h-[44px] z-30 flex items-center justify-between px-5 text-[8px] font-semibold text-slate-800 pointer-events-none bg-white">
            <span>9:41</span>
            {/* Dynamic Island */}
            <div className="w-[85px] h-[20px] bg-black rounded-full" />
            <div className="flex items-center gap-1">
              {/* Signal */}
              <div className="flex items-end gap-[1px] h-[8px]">
                <div className="w-[2px] h-[3px] bg-slate-800 rounded-[0.5px]" />
                <div className="w-[2px] h-[5px] bg-slate-800 rounded-[0.5px]" />
                <div className="w-[2px] h-[7px] bg-slate-800 rounded-[0.5px]" />
                <div className="w-[2px] h-[8px] bg-slate-800 rounded-[0.5px]" />
              </div>
              {/* Battery */}
              <div className="w-[18px] h-[9px] border border-slate-800 rounded-[2.5px] p-[1.5px] flex items-center relative ml-0.5">
                <div className="h-full w-[80%] bg-emerald-500 rounded-[1px]" />
                <div className="w-[1.5px] h-[4px] bg-slate-800/50 rounded-r-[1px] absolute -right-[2.5px]" />
              </div>
            </div>
          </div>

          {/* ── App Content ── */}
          <div className="flex-1 flex flex-col bg-white pt-[44px] pb-[58px] overflow-hidden">

            {/* App Header */}
            <div className="px-4 py-2 bg-white border-b border-slate-100/80 flex items-center justify-between shrink-0">
              <div className="flex items-center gap-1.5">
                <div className="w-[18px] h-[18px] bg-gradient-to-br from-orange-500 to-orange-600 rounded-md flex items-center justify-center">
                  <span className="text-[6px] font-black text-white leading-none">D</span>
                </div>
                <span className="text-[11px] font-extrabold tracking-tight text-slate-900">Devix</span>
              </div>
              <div className="flex items-center gap-2.5">
                <FiSearch className="w-3.5 h-3.5 text-slate-500" />
                <div className="relative">
                  <FiBell className="w-3.5 h-3.5 text-slate-500" />
                  <div className="absolute -top-0.5 -right-0.5 w-2 h-2 bg-orange-500 rounded-full border border-white flex items-center justify-center">
                    <span className="text-[4px] text-white font-bold">3</span>
                  </div>
                </div>
                <Avatar className="w-5 h-5 ring-2 ring-blue-500/30">
                  <AvatarImage src="https://api.dicebear.com/7.x/notionists/svg?seed=DevixUser" />
                </Avatar>
              </div>
            </div>

            {/* Feed Tabs */}
            <div className="px-4 py-2 bg-white flex items-center gap-4 text-[8px] font-semibold border-b border-slate-100/60 shrink-0">
              <span className="text-blue-600 border-b-[1.5px] border-blue-600 pb-1">Trending</span>
              <span className="text-slate-400 pb-1">Latest</span>
              <span className="text-slate-400 pb-1">Following</span>
              <div className="ml-auto flex items-center gap-1 text-[6.5px] text-emerald-600 font-bold bg-emerald-50 px-1.5 py-0.5 rounded-full">
                <div className="w-1 h-1 rounded-full bg-emerald-500 animate-pulse" />
                Live
              </div>
            </div>

            {/* Scrollable Feed */}
            <div className="flex-1 px-3 pt-2.5 pb-1 space-y-2 overflow-hidden bg-slate-50/70">

              {/* ── Post 1: Question with Code ── */}
              <motion.div
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.2 }}
                className="p-3 rounded-2xl bg-white border border-slate-100 shadow-[0_1px_3px_rgba(0,0,0,0.04)] space-y-2"
              >
                {/* Author */}
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-1.5">
                    <Avatar className="w-5 h-5 ring-1 ring-blue-100">
                      <AvatarImage src="https://api.dicebear.com/7.x/adventurer/svg?seed=Alex" />
                    </Avatar>
                    <div>
                      <div className="text-[8px] font-bold text-slate-800 flex items-center gap-0.5">
                        alex_dev
                        <div className="w-2.5 h-2.5 rounded-full bg-blue-500 text-white flex items-center justify-center">
                          <FiCheck className="w-1.5 h-1.5" />
                        </div>
                      </div>
                      <div className="text-[6px] text-slate-400 leading-none">Full Stack Developer · 2h</div>
                    </div>
                  </div>
                  <span className="text-[5.5px] px-1.5 py-0.5 rounded-full bg-blue-50 text-blue-600 font-bold border border-blue-100/50">Question</span>
                </div>

                {/* Post Content */}
                <p className="text-[8px] font-semibold text-slate-800 leading-[1.45]">
                  How to optimize WebSocket connections for real-time chat in production? 🔌
                </p>

                {/* Code Snippet */}
                <div className="p-2 rounded-xl bg-[#0d1117] font-mono text-[6.5px] leading-[1.6] text-slate-300 border border-slate-800/50 overflow-hidden">
                  <div className="flex items-center justify-between mb-1.5 pb-1 border-b border-slate-800/50">
                    <div className="flex items-center gap-1">
                      <div className="w-1.5 h-1.5 rounded-full bg-red-500/80" />
                      <div className="w-1.5 h-1.5 rounded-full bg-yellow-500/80" />
                      <div className="w-1.5 h-1.5 rounded-full bg-emerald-500/80" />
                    </div>
                    <span className="text-[5px] text-slate-500 font-sans">chat.ts</span>
                  </div>
                  <div><span className="text-purple-400">const</span> <span className="text-sky-300">ws</span> = <span className="text-purple-400">new</span> <span className="text-orange-300">WebSocket</span>(url);</div>
                  <div className="mt-0.5"><span className="text-sky-300">ws</span>.<span className="text-yellow-300">onmessage</span> = (<span className="text-orange-300">e</span>) {`=>`} &#123;</div>
                  <div className="pl-3"><span className="text-purple-400">const</span> msg = <span className="text-sky-300">JSON</span>.<span className="text-yellow-300">parse</span>(e.data);</div>
                  <div className="pl-3"><span className="text-yellow-300">dispatch</span>(<span className="text-emerald-400">addMessage</span>(msg));</div>
                  <div>&#125;;</div>
                </div>

                {/* Interaction Bar */}
                <div className="flex items-center justify-between pt-0.5">
                  <div className="flex items-center gap-3">
                    <span className="flex items-center gap-0.5 text-[7px] font-bold text-orange-500">
                      <FiArrowUp className="w-2.5 h-2.5" /> 48
                    </span>
                    <span className="flex items-center gap-0.5 text-[7px] text-slate-400">
                      <FiMessageCircle className="w-2.5 h-2.5" /> 12
                    </span>
                    <span className="flex items-center gap-0.5 text-[7px] text-slate-400">
                      <FiShare2 className="w-2.5 h-2.5" />
                    </span>
                  </div>
                  <FiBookmark className="w-2.5 h-2.5 text-slate-400" />
                </div>
              </motion.div>

              {/* ── Post 2: Build Log ── */}
              <motion.div
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.5 }}
                className="p-3 rounded-2xl bg-white border border-slate-100 shadow-[0_1px_3px_rgba(0,0,0,0.04)] space-y-1.5"
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-1.5">
                    <Avatar className="w-5 h-5 ring-1 ring-emerald-100">
                      <AvatarImage src="https://api.dicebear.com/7.x/adventurer/svg?seed=Sarah" />
                    </Avatar>
                    <div>
                      <div className="text-[8px] font-bold text-slate-800">sarah_builds</div>
                      <div className="text-[6px] text-slate-400 leading-none">DevOps Engineer · 5h</div>
                    </div>
                  </div>
                  <span className="text-[5.5px] px-1.5 py-0.5 rounded-full bg-emerald-50 text-emerald-600 font-bold border border-emerald-100/50">Build Log</span>
                </div>
                <p className="text-[8px] font-semibold text-slate-800 leading-[1.45]">
                  Day 47 of #100DaysOfCode — Built a real-time notification system using SSE! 🔔
                </p>
                <div className="flex items-center gap-1 flex-wrap">
                  <span className="text-[5.5px] px-1.5 py-0.5 rounded-full bg-slate-50 text-slate-500 font-semibold border border-slate-100">#devix</span>
                  <span className="text-[5.5px] px-1.5 py-0.5 rounded-full bg-slate-50 text-slate-500 font-semibold border border-slate-100">#sse</span>
                  <span className="text-[5.5px] px-1.5 py-0.5 rounded-full bg-slate-50 text-slate-500 font-semibold border border-slate-100">#react</span>
                </div>
                <div className="flex items-center justify-between pt-0.5">
                  <div className="flex items-center gap-3">
                    <span className="flex items-center gap-0.5 text-[7px] font-bold text-orange-500">
                      <FiArrowUp className="w-2.5 h-2.5" /> 87
                    </span>
                    <span className="flex items-center gap-0.5 text-[7px] text-slate-400">
                      <FiMessageCircle className="w-2.5 h-2.5" /> 23
                    </span>
                  </div>
                  <FiBookmark className="w-2.5 h-2.5 text-blue-500 fill-blue-500" />
                </div>
              </motion.div>

              {/* ── Live Chat Preview ── */}
              <motion.div
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.8 }}
                className="p-2.5 rounded-2xl bg-gradient-to-r from-blue-50 to-sky-50 border border-blue-100/60 space-y-1.5"
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-1">
                    <FiMessageCircle className="w-2.5 h-2.5 text-blue-600" />
                    <span className="text-[7px] font-bold text-blue-700">Direct Message</span>
                  </div>
                  <div className="flex items-center gap-0.5 text-[5.5px] text-emerald-600 font-bold bg-emerald-50 px-1 py-0.5 rounded-full">
                    <div className="w-1 h-1 rounded-full bg-emerald-500 animate-pulse" />
                    Online
                  </div>
                </div>
                <div className="flex items-start gap-1.5">
                  <Avatar className="w-4 h-4 shrink-0 ring-1 ring-blue-200">
                    <AvatarImage src="https://api.dicebear.com/7.x/adventurer/svg?seed=Emma" />
                  </Avatar>
                  <div className="py-1 px-2 rounded-xl rounded-tl-none bg-white border border-blue-100/80 text-[7px] text-slate-600 shadow-sm">
                    Hey! Loved your WebSocket approach. Can we pair on the auth module? 🚀
                  </div>
                </div>
                <div className="flex items-center gap-1 pl-5">
                  <div className="flex-1 py-1 px-2 rounded-full bg-white border border-blue-100/80 text-[6.5px] text-slate-400">Type a message...</div>
                  <div className="w-4 h-4 rounded-full bg-blue-600 flex items-center justify-center shrink-0 shadow-sm">
                    <FiSend className="w-2 h-2 text-white" />
                  </div>
                </div>
              </motion.div>

            </div>
          </div>

          {/* ── Bottom Tab Bar ── */}
          <div className="absolute bottom-0 inset-x-0 bg-white border-t border-slate-100 px-2 pt-1.5 pb-5 flex items-center justify-around z-20">
            <div className="flex flex-col items-center gap-[2px]">
              <FiHome className="w-3.5 h-3.5 text-blue-600" />
              <span className="text-[5px] font-bold text-blue-600">Home</span>
            </div>
            <div className="flex flex-col items-center gap-[2px]">
              <FiCompass className="w-3.5 h-3.5 text-slate-400" />
              <span className="text-[5px] font-medium text-slate-400">Explore</span>
            </div>
            <div className="flex flex-col items-center gap-[2px] -mt-3.5">
              <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-blue-600 to-sky-500 flex items-center justify-center shadow-lg shadow-blue-500/30 ring-2 ring-white">
                <FiPlusCircle className="w-4 h-4 text-white" />
              </div>
              <span className="text-[5px] font-medium text-slate-400">Post</span>
            </div>
            <div className="flex flex-col items-center gap-[2px] relative">
              <FiMessageCircle className="w-3.5 h-3.5 text-slate-400" />
              <span className="text-[5px] font-medium text-slate-400">Chat</span>
              <div className="absolute -top-0.5 right-0.5 w-1.5 h-1.5 bg-orange-500 rounded-full border border-white" />
            </div>
            <div className="flex flex-col items-center gap-[2px]">
              <FiUser className="w-3.5 h-3.5 text-slate-400" />
              <span className="text-[5px] font-medium text-slate-400">Profile</span>
            </div>
          </div>

          {/* Home Indicator */}
          <div className="absolute bottom-1.5 left-1/2 -translate-x-1/2 w-24 h-1 bg-slate-200 rounded-full z-30" />
        </motion.div>

        {/* ────── Floating Widgets ────── */}

        {/* Widget 1: Live notification */}
        <motion.a
          href="https://github.com/neutron420/devix-frontend"
          target="_blank"
          rel="noopener noreferrer"
          animate={{ y: [0, -7, 0], x: [0, 3, 0] }}
          transition={{ repeat: Infinity, duration: 5.5, ease: "easeInOut" }}
          className="absolute top-[3%] -left-20 p-2.5 rounded-2xl bg-white/95 border border-white/30 backdrop-blur-md shadow-[0_8px_30px_rgba(0,0,0,0.12)] flex items-center gap-2.5 w-[190px] z-10 hover:scale-105 active:scale-95 transition-transform cursor-pointer"
        >
          <Avatar className="w-6.5 h-6.5 border border-orange-200 shrink-0">
            <AvatarImage src="https://api.dicebear.com/7.x/adventurer/svg?seed=Rahul" />
          </Avatar>
          <div className="min-w-0 text-left">
            <p className="text-[9.5px] font-bold text-slate-800 flex items-center gap-1">
              New Post <FiZap className="w-2.5 h-2.5 text-orange-500 animate-pulse" />
            </p>
            <p className="text-[8.5px] text-slate-500 truncate font-medium">Rahul shared a build log</p>
          </div>
        </motion.a>

        {/* Widget 2: Next.js tag */}
        <motion.a
          href="https://nextjs.org"
          target="_blank"
          rel="noopener noreferrer"
          animate={{ y: [0, 8, 0], x: [0, -3, 0] }}
          transition={{ repeat: Infinity, duration: 6, ease: "easeInOut", delay: 0.5 }}
          className="absolute top-[18%] -right-16 py-1.5 px-3 rounded-full bg-white/95 border border-white/30 backdrop-blur-md shadow-[0_8px_30px_rgba(0,0,0,0.12)] flex items-center gap-1.5 text-[9px] font-bold text-slate-800 z-10 hover:scale-105 active:scale-95 transition-transform cursor-pointer"
        >
          <div className="w-4 h-4 rounded-full bg-slate-900 text-white flex items-center justify-center shrink-0">
            <SiNextdotjs className="w-2.5 h-2.5" />
          </div>
          <span>Next.js 16</span>
        </motion.a>

        {/* Widget 3: Tailwind tag */}
        <motion.a
          href="https://tailwindcss.com"
          target="_blank"
          rel="noopener noreferrer"
          animate={{ y: [0, -6, 0], x: [0, 4, 0] }}
          transition={{ repeat: Infinity, duration: 5.2, ease: "easeInOut", delay: 1.2 }}
          className="absolute top-[32%] -right-12 py-1.5 px-3 rounded-full bg-white/95 border border-white/30 backdrop-blur-md shadow-[0_8px_30px_rgba(0,0,0,0.12)] flex items-center gap-1.5 text-[9px] font-bold text-slate-800 z-10 hover:scale-105 active:scale-95 transition-transform cursor-pointer"
        >
          <SiTailwindcss className="w-3.5 h-3.5 text-sky-400" />
          <span>Tailwind v4</span>
        </motion.a>

        {/* Widget 4: Upvotes */}
        <motion.div
          animate={{ y: [0, -9, 0] }}
          transition={{ repeat: Infinity, duration: 4.8, ease: "easeInOut", delay: 1 }}
          className="absolute bottom-[28%] -left-18 py-1.5 px-3.5 rounded-full bg-white/95 border border-white/30 backdrop-blur-md shadow-[0_8px_30px_rgba(0,0,0,0.12)] flex items-center gap-2 text-[10px] font-bold text-orange-500 z-10"
        >
          <FiArrowUp className="w-3.5 h-3.5 animate-bounce" style={{ animationDuration: "1.5s" }} />
          <span>+142 Upvotes</span>
        </motion.div>

        {/* Widget 5: Code editor */}
        <motion.div
          animate={{ y: [0, 7, 0] }}
          transition={{ repeat: Infinity, duration: 6.2, ease: "easeInOut", delay: 1.8 }}
          className="absolute bottom-[24%] -right-16 py-2 px-3 rounded-2xl bg-white/95 border border-white/30 backdrop-blur-md shadow-[0_8px_30px_rgba(0,0,0,0.12)] flex items-center gap-2 w-[170px] z-10"
        >
          <div className="p-1.5 rounded-lg bg-blue-50 text-blue-500 shrink-0">
            <FiCode className="w-3.5 h-3.5" />
          </div>
          <div className="min-w-0 text-left">
            <p className="text-[9.5px] font-bold text-slate-800 leading-none">Markdown Editor</p>
            <p className="text-[8px] text-slate-400 truncate mt-1">Write & share code articles</p>
          </div>
        </motion.div>

        {/* Widget 6: Active devs */}
        <motion.div
          animate={{ y: [0, -5, 0], x: [0, -4, 0] }}
          transition={{ repeat: Infinity, duration: 5.8, ease: "easeInOut", delay: 0.8 }}
          className="absolute bottom-[8%] -left-14 p-2 rounded-2xl bg-white/95 border border-white/30 backdrop-blur-md shadow-[0_8px_30px_rgba(0,0,0,0.12)] flex items-center gap-2 w-[160px] z-10"
        >
          <div className="p-1.5 rounded-lg bg-emerald-50 text-emerald-600 shrink-0">
            <FiUsers className="w-3.5 h-3.5" />
          </div>
          <div className="text-left">
            <p className="text-[9px] font-bold text-slate-800 leading-none">Devix Community</p>
            <p className="text-[8px] text-slate-400 mt-1 flex items-center gap-1 font-semibold">
              <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-ping" />
              14 Online
            </p>
          </div>
        </motion.div>

        {/* Floating emojis */}
        <motion.div
          animate={{ y: [0, 10, 0] }}
          transition={{ repeat: Infinity, duration: 4.2, ease: "easeInOut" }}
          className="absolute top-[48%] -left-12 w-7 h-7 rounded-full bg-white/90 shadow-md flex items-center justify-center text-[11px] z-10"
        >
          🚀
        </motion.div>
        <motion.div
          animate={{ y: [0, -8, 0] }}
          transition={{ repeat: Infinity, duration: 3.8, ease: "easeInOut", delay: 0.5 }}
          className="absolute bottom-[12%] right-24 w-7 h-7 rounded-full bg-white/90 shadow-md flex items-center justify-center text-[11px] z-10"
        >
          💬
        </motion.div>
      </div>

      {/* ── App Store & Play Store Badges ── */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 1.2, duration: 0.6 }}
        className="relative z-10 flex flex-col items-center gap-3 pb-8 mt-auto"
      >
        <p className="text-[11px] font-semibold text-white/80 tracking-wide uppercase">Coming Soon on</p>
        <div className="flex items-center gap-3">
          {/* App Store Badge */}
          <a
            href="#"
            className="flex items-center gap-2 bg-black/80 backdrop-blur-md rounded-xl px-3.5 py-2 border border-white/10 hover:bg-black/90 hover:border-white/20 transition-all hover:scale-105 active:scale-95 cursor-pointer shadow-lg"
          >
            <svg className="w-5 h-5 text-white" viewBox="0 0 24 24" fill="currentColor">
              <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.8-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>
            </svg>
            <div className="text-left">
              <p className="text-[5.5px] text-white/60 font-medium leading-none uppercase tracking-wider">Download on the</p>
              <p className="text-[10px] text-white font-bold leading-tight">App Store</p>
            </div>
          </a>
          {/* Google Play Badge */}
          <a
            href="#"
            className="flex items-center gap-2 bg-black/80 backdrop-blur-md rounded-xl px-3.5 py-2 border border-white/10 hover:bg-black/90 hover:border-white/20 transition-all hover:scale-105 active:scale-95 cursor-pointer shadow-lg"
          >
            <svg className="w-5 h-5" viewBox="0 0 24 24" fill="none">
              <path d="M3.609 1.814L13.793 12l-10.184 10.186a.916.916 0 0 1-.609-.862V2.676c0-.335.184-.635.609-.862z" fill="#4285F4"/>
              <path d="M17.217 8.577l-3.424 3.424 3.424 3.424 3.862-2.174c.44-.247.704-.624.704-1.25s-.264-1.003-.704-1.25l-3.862-2.174z" fill="#FBBC04"/>
              <path d="M3.609 1.814L13.793 12l3.424-3.424L6.216.576c-.462-.26-.967-.26-1.463 0l-1.144.644.609.594z" fill="#EA4335"/>
              <path d="M3.609 22.186L13.793 12l3.424 3.424L6.216 23.424c-.462.26-.967.26-1.463 0l-1.144-.644z" fill="#34A853"/>
            </svg>
            <div className="text-left">
              <p className="text-[5.5px] text-white/60 font-medium leading-none uppercase tracking-wider">Get it on</p>
              <p className="text-[10px] text-white font-bold leading-tight">Google Play</p>
            </div>
          </a>
        </div>
      </motion.div>
    </div>
  );
}
