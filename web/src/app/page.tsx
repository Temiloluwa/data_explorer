import MaxWidthWrapper from "@/components/MaxWidthWrapper";
import { buttonVariants } from "@/components/ui/button";
import { ArrowRight } from "lucide-react";
import Image from "next/image";
import Link from "next/link";

export default function Home() {


  return (
    <div>
      <>
        <MaxWidthWrapper className="mb-12 mt-28 sm:mt-40 flex flex-col items-center justify-center text-center">
          <div className="mx-auto mb-4 flex max-w-fit items-center justify-center space-x-2 overflow-hidden rounded-full border-gray-200 bg-white px-7 py-2 shadow-md backdrop-blur transition-all hover:border-gray-300 hover:bg-white/50">
            <p className="text-sm font-semibold text-gray-700">Data Explorer is Live!</p>
          </div>
          <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold text-center text-gray-900">
            <span className="block">Unlock the Power of Your Data 🚀</span>
            <span className="block text-purple-600">Discover Insights, Drive Innovation 💡</span>
          </h1>

          <div className="max-w-prose mt-5 sm:text-large text-zinc-700">
            <p>Our AI-powered app enables seamless data interaction, providing automatic insights and trend alerts in real-time. Whether you're an analyst, scientist, or decision-maker, our platform streamlines analysis, delivering actionable insights with ease.</p>
          </div>

          <Link className={buttonVariants({size: "lg", className: "mt-5"})} href="/dashboard"  target="_blank">
            Discover More <ArrowRight className="ml-2 h-5 w-5"/>
          </Link>
        </MaxWidthWrapper>

        <div className="relative isolate flex items-center justify-center">
          <div className="mx-auto max-w-6x1 px-6 lg:px-8">
            <div className="mt-16 flow-root sm:mt-24">
              <div className="-m-2 rounded-xl bg-gray-900/5 p-2 ring-1 ring-inset ring-gray-900/10 lg:-m-4 lg:rounded-2xl lg:p-4"> 
                <Image 
                  alt="landing_page_image" 
                  src="/landing_page_bg_3.jpg" 
                  height={500} 
                  width={900} 
                  quality={100} 
                  className="mx-auto" 
                />
              </div>
            </div>
          </div>
        </div>

        <div className='mx-auto mb-32 mt-32 max-w-5xl sm:mt-56'>
          <div className='mb-12 px-6 lg:px-8'>
            <div className='mx-auto max-w-2xl sm:text-center'>
              <h2 className='mt-2 font-bold text-4xl text-gray-900 sm:text-5xl'>
                Explore with Data Explorer
              </h2>
              <p className='mt-4 text-lg text-gray-600'>
                Uncover hidden insights and trends in your data effortlessly with Data Explorer.
              </p>
            </div>
          </div>

          {/* steps */}
          <ol className='my-8 md:flex md:space-x-12'>
            <li className='md:flex-1'>
              <div className='flex flex-col space-y-2 border-l-4 border-zinc-300 py-2 pl-4 md:border-l-0 md:border-t-2 md:pb-0 md:pl-0 md:pt-4'>
                <span className='text-sm font-medium text-blue-600'>
                  Step 1
                </span>
                <span className='text-xl font-semibold'>
                  Sign up for an account
                </span>
                <span className='mt-2 text-zinc-700'>
                  Get started with Data Explorer and unlock the power of your data today.
                </span>
              </div>
            </li>
            <li className='md:flex-1'>
              <div className='flex flex-col space-y-2 border-l-4 border-zinc-300 py-2 pl-4 md:border-l-0 md:border-t-2 md:pb-0 md:pl-0 md:pt-4'>
                <span className='text-sm font-medium text-blue-600'>
                  Step 2
                </span>
                <span className='text-xl font-semibold'>
                  Upload your data
                </span>
                <span className='mt-2 text-zinc-700'>
                  Easily upload your datasets and let Data Explorer do the rest.
                </span>
              </div>
            </li>
            <li className='md:flex-1'>
              <div className='flex flex-col space-y-2 border-l-4 border-zinc-300 py-2 pl-4 md:border-l-0 md:border-t-2 md:pb-0 md:pl-0 md:pt-4'>
                <span className='text-sm font-medium text-blue-600'>
                  Step 3
                </span>
                <span className='text-xl font-semibold'>
                  Explore your data
                </span>
                <span className='mt-2 text-zinc-700'>
                  Dive into your data like never before and discover insights with Data Explorer.
                </span>
              </div>
            </li>
          </ol>
        </div>

        <footer className="bg-gray-900 text-white py-8">
          <div className="container mx-auto px-4">
            <div className="flex flex-col md:flex-row items-center justify-between">
              <p className="text-lg">© 2024 Data Explorer. All rights reserved.</p>
              <ul className="flex space-x-4 mt-4 md:mt-0">
                <li>
                  <a href="#" className="hover:text-gray-400">Terms of Service</a>
                </li>
                <li>
                  <a href="#" className="hover:text-gray-400">Privacy Policy</a>
                </li>
                <li>
                  <a href="#" className="hover:text-gray-400">Contact Us</a>
                </li>
              </ul>
            </div>
          </div>
        </footer>

      </>
    </div>
  );
}
