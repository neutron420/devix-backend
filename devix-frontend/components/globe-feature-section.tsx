"use client";

import { Button } from "@/components/ui/button";
import { ArrowRight } from "lucide-react";
import { Globe } from "@/components/ui/cobe-globe";

const markers = [
  { id: "sf", location: [37.7595, -122.4367] as [number, number], label: "San Francisco" },
  { id: "nyc", location: [40.7128, -74.006] as [number, number], label: "New York" },
  { id: "tokyo", location: [35.6762, 139.6503] as [number, number], label: "Tokyo" },
  { id: "london", location: [51.5074, -0.1278] as [number, number], label: "London" },
  { id: "sydney", location: [-33.8688, 151.2093] as [number, number], label: "Sydney" },
  { id: "capetown", location: [-33.9249, 18.4241] as [number, number], label: "Cape Town" },
  { id: "dubai", location: [25.2048, 55.2708] as [number, number], label: "Dubai" },
  { id: "saopaulo", location: [-23.5505, -46.6333] as [number, number], label: "São Paulo" },
]

const arcs = [
  {
    id: "sf-tokyo",
    from: [37.7595, -122.4367] as [number, number],
    to: [35.6762, 139.6503] as [number, number],
    label: "SF → Tokyo",
  },
  {
    id: "nyc-london",
    from: [40.7128, -74.006] as [number, number],
    to: [51.5074, -0.1278] as [number, number],
    label: "NYC → London",
  },
]

export default function Featured_05() {
  return (
    <section className="bg-white dark:bg-black dark:bg-transparent border-b border-gray-200 dark:border-gray-800">
      <div className="mx-auto container border-x border-gray-200 dark:border-gray-800 py-6 px-4">
        <div className="relative border border-dashed border-zinc-400 dark:border-zinc-700 rounded-lg p-6 md:p-16 bg-white dark:bg-zinc-950 text-black dark:text-white">
          <CornerPlusIcons />
          <div className="flex flex-col-reverse items-center justify-between gap-10 md:flex-row relative z-10">
            <div className="z-10 max-w-xl text-center md:text-left flex flex-col items-center md:items-start">
              <h1 className="text-3xl md:text-5xl font-black uppercase tracking-tighter text-black dark:text-white mb-4 w-full" style={{ fontFamily: '"Arial Black", Impact, sans-serif' }}>
                CONNECTING BUILDERS <span className="text-[#fd5200] italic">GLOBALLY</span>
              </h1>
              <p className="text-zinc-700 dark:text-zinc-400 text-sm md:text-base leading-relaxed mb-6">
                Devix bridges geographic boundaries to connect engineers, creators, and open-source contributors. Collaborate in real-time, share your daily progress, and build your reputation across a decentralized global developer network.
              </p>
              <Button className="inline-flex items-center gap-2 rounded-full bg-[#0038FF] text-white hover:bg-blue-700 px-5 py-2 text-sm font-bold shadow-md transition duration-300">
                Start Building <ArrowRight className="h-4 w-4" />
              </Button>
            </div>
            <div className="relative h-[300px] md:h-[350px] w-full max-w-md flex items-center justify-center">
              <Globe
                className="w-[280px] h-[280px] md:w-[350px] md:h-[350px]"
                markers={markers}
                arcs={arcs}
                markerColor={[253 / 255, 82 / 255, 0]} // Devix Orange
                baseColor={[1, 1, 1]}
                arcColor={[0, 56 / 255, 1]} // Devix Blue
                glowColor={[0.94, 0.93, 0.91]}
                dark={0}
                mapBrightness={10}
                markerSize={0.025}
                markerElevation={0.01}
              />
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

const CornerPlusIcons = () => (
  <>
    <PlusIcon className="absolute -top-3 -left-3" />
    <PlusIcon className="absolute -top-3 -right-3" />
    <PlusIcon className="absolute -bottom-3 -left-3" />
    <PlusIcon className="absolute -bottom-3 -right-3" />
  </>
)

const PlusIcon = ({ className }: { className?: string }) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    fill="none"
    viewBox="0 0 24 24"
    width={24}
    height={24}
    strokeWidth="1"
    stroke="currentColor"
    className={`dark:text-white text-black size-6 ${className}`}
  >
    <path strokeLinecap="round" strokeLinejoin="round" d="M12 6v12m6-6H6" />
  </svg>
)
