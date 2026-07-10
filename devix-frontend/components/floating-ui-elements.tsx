"use client";

import { motion } from "framer-motion";
import {
  FiMessageCircle,
  FiHeart,
  FiShare2,
  FiAward,
  FiZap,
  FiUsers,
  FiHash,
} from "react-icons/fi";
import { SiNextdotjs, SiGo } from "react-icons/si";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

/* ── premium glassmorphic styling system ── */
const glass =
  "bg-white/40 dark:bg-neutral-950/45 backdrop-blur-2xl border border-white/25 dark:border-white/10 shadow-[0_8px_32px_rgba(0,0,0,0.06)] dark:shadow-[0_8px_32px_rgba(0,0,0,0.45)] transition-all duration-300 hover:border-white/50 dark:hover:border-white/20 hover:shadow-[0_12px_48px_rgba(0,0,0,0.12)] dark:hover:shadow-[0_12px_48px_rgba(0,0,0,0.7)] pointer-events-auto";

const pill =
  "bg-white/50 dark:bg-neutral-950/50 backdrop-blur-2xl border border-white/20 dark:border-white/10 shadow-sm transition-all duration-300 hover:border-white/40 dark:hover:border-white/20 hover:scale-105 pointer-events-auto";

export function FloatingUIElements() {
  return (
    <div className="absolute inset-0 z-0 pointer-events-none overflow-hidden select-none">

      {/* ───────── LEFT GUTTER (Vertical Alignment) ───────── */}


      {/* 2. Trending Tag - #system-design (Left, top-[30%]) */}
      <motion.div
        initial={{ opacity: 0, x: -15 }}
        animate={{ opacity: 1, x: 0, y: [4, -4] }}
        transition={{ repeat: Infinity, repeatType: "mirror", duration: 6, ease: "easeInOut", delay: 0.8 }}
        className={`absolute top-[30%] left-4 lg:left-8 xl:left-12 2xl:left-16 py-1.5 px-3 rounded-full ${pill} hidden lg:flex items-center gap-1 text-[10px] font-semibold text-blue-600 dark:text-blue-400`}
      >
        <FiHash className="w-3 h-3 text-blue-500" />
        <span>system-design</span>
      </motion.div>

      {/* 3. Achievement Badge (Left, top-[42%]) */}
      <motion.div
        initial={{ opacity: 0, scale: 0.9 }}
        animate={{ opacity: 1, scale: 1, y: [-5, 5] }}
        transition={{ repeat: Infinity, repeatType: "mirror", duration: 4.5, ease: "easeInOut", delay: 0.5 }}
        className="absolute top-[42%] left-4 lg:left-8 xl:left-12 2xl:left-16 p-3 rounded-xl bg-gradient-to-br from-amber-50/60 to-yellow-50/60 dark:from-amber-950/30 dark:to-yellow-900/30 backdrop-blur-2xl border border-amber-200/25 dark:border-amber-700/25 shadow-[0_4px_20px_rgba(245,158,11,0.08)] hidden md:block pointer-events-auto hover:scale-105 transition-transform duration-300 origin-left"
      >
        <div className="flex items-center gap-2">
          <FiAward className="w-4 h-4 text-amber-500" />
          <div>
            <p className="text-[10px] font-bold text-amber-700 dark:text-amber-300">Top Contributor</p>
            <p className="text-[8px] text-amber-600/70 dark:text-amber-400/60 mt-0.5">🏆 Bug Hunter</p>
          </div>
        </div>
      </motion.div>

      {/* 4. Glowing Avatar (Left, top-[56%]) */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1, y: [8, -8], rotate: [2, -2] }}
        transition={{ repeat: Infinity, repeatType: "mirror", duration: 5.5, ease: "easeInOut", delay: 1.2 }}
        className="absolute top-[56%] left-4 lg:left-8 xl:left-12 2xl:left-20 hidden sm:block pointer-events-auto"
      >
        <div className="relative group cursor-pointer">
          <div className="absolute -inset-1 rounded-full bg-gradient-to-r from-orange-500 to-amber-500 blur-md opacity-25 group-hover:opacity-40 transition-opacity duration-300" />
          <Avatar className="w-11 h-11 border-2 border-orange-400/40 shadow-lg relative z-10 bg-white dark:bg-neutral-900 group-hover:scale-105 transition-transform duration-300">
            <AvatarImage src="https://api.dicebear.com/7.x/notionists/svg?seed=Devix" />
            <AvatarFallback>DV</AvatarFallback>
          </Avatar>
        </div>
      </motion.div>

      {/* 5. Comment Card 2 (Left, top-[70%]) */}
      <motion.div
        initial={{ opacity: 0, x: -20 }}
        animate={{ opacity: 1, x: 0, y: [-4, 4] }}
        transition={{ repeat: Infinity, repeatType: "mirror", duration: 6.5, ease: "easeInOut", delay: 0.6 }}
        className={`absolute top-[70%] left-4 lg:left-8 xl:left-12 2xl:left-16 p-3 rounded-xl ${glass} w-[200px] hidden 2xl:block`}
      >
        <div className="flex items-start gap-2">
          <Avatar className="w-6 h-6 border border-blue-400/30 shrink-0">
            <AvatarImage src="https://api.dicebear.com/7.x/adventurer/svg?seed=Felix" />
            <AvatarFallback>FX</AvatarFallback>
          </Avatar>
          <div className="min-w-0">
            <p className="text-[10px] font-bold text-foreground truncate">Just pushed a fix!</p>
            <p className="text-[9px] text-muted-foreground mt-0.5 leading-relaxed">Auth bug resolved in #42</p>
          </div>
        </div>
        <div className="flex items-center gap-3 mt-2 pl-8 text-[8px] text-muted-foreground font-semibold">
          <span className="flex items-center gap-0.5 text-rose-500/90"><FiHeart className="w-2.5 h-2.5" /> 12</span>
          <span className="flex items-center gap-0.5 text-sky-500/90"><FiShare2 className="w-2.5 h-2.5" /> Share</span>
        </div>
      </motion.div>

      {/* 6. Typing Indicator (Left, bottom-[8%]) */}
      <motion.div
        initial={{ opacity: 0, y: 15 }}
        animate={{ opacity: 1, y: 0 }}
        className={`absolute bottom-[8%] left-4 lg:left-8 xl:left-12 2xl:left-16 py-2 px-3.5 rounded-full ${pill} hidden sm:flex items-center gap-2`}
      >
        <Avatar className="w-4 h-4">
          <AvatarImage src="https://api.dicebear.com/7.x/notionists/svg?seed=Sarah" />
        </Avatar>
        <span className="text-[10px] font-semibold text-muted-foreground">Sarah is typing</span>
        <span className="flex gap-[2px]">
          <span className="w-0.75 h-0.75 rounded-full bg-muted-foreground/60 animate-bounce [animation-delay:0ms]" />
          <span className="w-0.75 h-0.75 rounded-full bg-muted-foreground/60 animate-bounce [animation-delay:150ms]" />
          <span className="w-0.75 h-0.75 rounded-full bg-muted-foreground/60 animate-bounce [animation-delay:300ms]" />
        </span>
      </motion.div>


      {/* ───────── RIGHT GUTTER (Vertical Alignment) ───────── */}

      {/* 7. Comment Card 1 (Right, top-[18%]) */}
      <motion.div
        initial={{ opacity: 0, x: 20 }}
        animate={{ opacity: 1, x: 0, y: [-4, 4] }}
        transition={{ repeat: Infinity, repeatType: "mirror", duration: 5.2, ease: "easeInOut", delay: 0.1 }}
        className={`absolute top-[18%] right-4 lg:right-8 xl:right-12 2xl:right-16 p-3 rounded-xl ${glass} w-[220px] hidden lg:block origin-right scale-90 xl:scale-95 2xl:scale-100`}
      >
        <div className="flex items-start gap-2">
          <Avatar className="w-6.5 h-6.5 border border-orange-400/30 shrink-0">
            <AvatarImage src="https://api.dicebear.com/7.x/adventurer/svg?seed=Rahul" />
            <AvatarFallback>RK</AvatarFallback>
          </Avatar>
          <div className="min-w-0">
            <p className="text-[10px] font-bold text-foreground flex items-center gap-1 truncate">
              New API deployed <FiZap className="w-2.5 h-2.5 text-orange-500 shrink-0" />
            </p>
            <p className="text-[9px] text-muted-foreground mt-0.5 leading-normal">Speed is up 40%! ⚡</p>
          </div>
        </div>
        <div className="flex items-center gap-3 mt-2 pl-8 text-[8px] text-muted-foreground font-semibold">
          <span className="flex items-center gap-0.5 text-rose-500/90"><FiHeart className="w-2.5 h-2.5" /> 24</span>
          <span className="flex items-center gap-0.5 text-sky-500/90"><FiMessageCircle className="w-2.5 h-2.5" /> 5</span>
        </div>
      </motion.div>

      {/* 8. Trending Tag - #100DaysOfCode (Right, top-[30%]) */}
      <motion.div
        initial={{ opacity: 0, x: 15 }}
        animate={{ opacity: 1, x: 0, y: [-3, 3] }}
        transition={{ repeat: Infinity, repeatType: "mirror", duration: 5.8, ease: "easeInOut", delay: 0.9 }}
        className={`absolute top-[30%] right-4 lg:right-8 xl:right-12 2xl:right-16 py-1.5 px-3 rounded-full ${pill} hidden md:flex items-center gap-1 text-[10px] font-semibold text-orange-600 dark:text-orange-400`}
      >
        <FiHash className="w-3 h-3 text-orange-500" />
        <span>100DaysOfCode</span>
      </motion.div>

      {/* 9. Mini Poll (Right, top-[42%]) */}
      <motion.div
        initial={{ opacity: 0, y: 15 }}
        animate={{ opacity: 1, y: [3, -3] }}
        transition={{ repeat: Infinity, repeatType: "mirror", duration: 6.2, ease: "easeInOut", delay: 0.4 }}
        className={`absolute top-[42%] right-4 lg:right-8 xl:right-12 2xl:right-16 p-3.5 rounded-xl ${glass} w-[190px] hidden xl:block origin-right scale-90 xl:scale-95 2xl:scale-100`}
      >
        <div className="flex items-center gap-1.5 mb-2">
          <p className="text-[10px] font-bold text-foreground">Next.js vs Go?</p>
        </div>
        
        <div className="space-y-2">
          {/* Next.js Option */}
          <div>
            <div className="flex justify-between items-center text-[9px] text-muted-foreground mb-0.5 font-medium">
              <div className="flex items-center gap-1 font-semibold text-foreground">
                <SiNextdotjs className="w-3 h-3 text-black dark:text-white shrink-0" />
                <span>Next.js</span>
              </div>
              <span className="font-bold text-blue-500">62%</span>
            </div>
            <div className="h-1 bg-neutral-200 dark:bg-neutral-800 rounded-full overflow-hidden">
              <motion.div
                initial={{ width: 0 }}
                animate={{ width: "62%" }}
                transition={{ duration: 1.2, ease: "easeOut", delay: 1 }}
                className="h-full bg-blue-500 rounded-full"
              />
            </div>
          </div>

          {/* Go Option */}
          <div>
            <div className="flex justify-between items-center text-[9px] text-muted-foreground mb-0.5 font-medium">
              <div className="flex items-center gap-1 font-semibold text-foreground">
                <SiGo className="w-3.5 h-3.5 text-cyan-500 shrink-0" />
                <span>GoLang</span>
              </div>
              <span className="font-bold text-emerald-500">38%</span>
            </div>
            <div className="h-1 bg-neutral-200 dark:bg-neutral-800 rounded-full overflow-hidden">
              <motion.div
                initial={{ width: 0 }}
                animate={{ width: "38%" }}
                transition={{ duration: 1.2, ease: "easeOut", delay: 1.2 }}
                className="h-full bg-cyan-500 rounded-full"
              />
            </div>
          </div>
        </div>
        <p className="text-[7.5px] text-muted-foreground/60 mt-1.5 pl-0.5 font-medium">247 active votes</p>
      </motion.div>

      {/* 10. Organization Card (Right, top-[62%]) */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1, y: [-5, 5] }}
        transition={{ repeat: Infinity, repeatType: "mirror", duration: 7, ease: "easeInOut", delay: 1.5 }}
        className={`absolute top-[62%] right-4 lg:right-8 xl:right-12 2xl:right-16 p-3 rounded-xl ${glass} w-[190px] hidden 2xl:block`}
      >
        <div className="flex items-center gap-2">
          <div className="p-1.5 rounded-lg bg-violet-500/10">
            <FiUsers className="w-3.5 h-3.5 text-violet-500" />
          </div>
          <div>
            <p className="text-[11px] font-bold text-foreground leading-tight">Devix Core Team</p>
            <p className="text-[8px] text-muted-foreground mt-0.5">14 devs online</p>
          </div>
        </div>
        <div className="flex -space-x-1.5 mt-2.5 pl-0.5">
          <Avatar className="w-5.5 h-5.5 border border-white dark:border-neutral-900"><AvatarImage src="https://api.dicebear.com/7.x/adventurer/svg?seed=Anna" /></Avatar>
          <Avatar className="w-5.5 h-5.5 border border-white dark:border-neutral-900"><AvatarImage src="https://api.dicebear.com/7.x/adventurer/svg?seed=Bob" /></Avatar>
          <Avatar className="w-5.5 h-5.5 border border-white dark:border-neutral-900"><AvatarImage src="https://api.dicebear.com/7.x/adventurer/svg?seed=Charlie" /></Avatar>
          <Avatar className="w-5.5 h-5.5 border border-white dark:border-neutral-900"><AvatarImage src="https://api.dicebear.com/7.x/adventurer/svg?seed=Mia" /></Avatar>
          <div className="w-5.5 h-5.5 rounded-full border border-white dark:border-neutral-900 bg-neutral-100 dark:bg-neutral-800 flex items-center justify-center text-[7px] font-bold text-muted-foreground">+10</div>
        </div>
      </motion.div>

      {/* 11. Heart Badge (Right, top-[78%]) */}
      <motion.div
        initial={{ opacity: 0, scale: 0.85 }}
        animate={{ opacity: 1, scale: 1, y: [-4, 4] }}
        transition={{ repeat: Infinity, repeatType: "mirror", duration: 4.8, ease: "easeInOut", delay: 1.8 }}
        className={`absolute top-[78%] right-4 lg:right-8 xl:right-12 2xl:right-24 py-1.5 px-3 rounded-full ${pill} hidden sm:flex items-center gap-1.5`}
      >
        <FiHeart className="w-3 h-3 text-rose-500 animate-pulse" />
        <span className="text-[10px] font-bold text-foreground">128</span>
      </motion.div>



      {/* ───────── HEADER WIDGET (Top Right) ───────── */}


    </div>
  );
}
