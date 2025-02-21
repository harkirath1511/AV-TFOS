import React from 'react';

function Footer() {
  return (
    <footer className="w-full bg-gradient-to-r from-gray-800 via-gray-900 to-black text-white py-8">
      <div className="max-w-7xl mx-auto px-6">
        
        {/* Top Section */}
        <div className="flex flex-col lg:flex-row justify-between items-center text-center lg:text-left">
          
          {/* Left - Branding & Socials */}
          <div className="mb-6 lg:mb-0">
            <h4 className="text-2xl font-bold">Let's keep in touch!</h4>
            <p className="text-gray-400">Find us on these platforms:</p>
            <div className="flex mt-3 space-x-4">
              <button onClick={() => window.open('https://github.com/vaidikcode')} className="bg-white p-2 rounded-full shadow-md">
                <img src="https://cdn-icons-png.flaticon.com/512/25/25231.png" alt="GitHub" className="w-7 h-7"/>
              </button>
              <button onClick={() => window.open('https://slack.com')} className="bg-white p-2 rounded-full shadow-md">
                <img src="https://cdn-icons-png.flaticon.com/512/2111/2111615.png" alt="Slack" className="w-7 h-7"/>
              </button>
              <button onClick={() => window.open('https://linkedin.com')} className="bg-white p-2 rounded-full shadow-md">
                <img src="https://cdn-icons-png.flaticon.com/512/145/145807.png" alt="LinkedIn" className="w-7 h-7"/>
              </button>
              <button className="bg-white p-2 rounded-full shadow-md">
                <img src="https://pbs.twimg.com/profile_images/1683899100922511378/5lY42eHs_400x400.jpg" alt="Other" className="w-7 h-7"/>
              </button>
            </div>
          </div>

          {/* Right - Links */}
          <div className="grid grid-cols-2 gap-10">
            <div>
              <h5 className="font-semibold text-gray-300">Useful Links</h5>
              <ul className="mt-2 space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-white">About Us</a></li>
                <li><a href="#" className="hover:text-white">Blog</a></li>
                <li><a href="#" className="hover:text-white">GitHub</a></li>
                <li><a href="#" className="hover:text-white">Free Products</a></li>
              </ul>
            </div>
            <div>
              <h5 className="font-semibold text-gray-300">Other Resources</h5>
              <ul className="mt-2 space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-white">MIT License</a></li>
                <li><a href="#" className="hover:text-white">Terms & Conditions</a></li>
                <li><a href="#" className="hover:text-white">Privacy Policy</a></li>
                <li><a href="#" className="hover:text-white">Contact Us</a></li>
              </ul>
            </div>
          </div>
        </div>

        {/* Divider */}
        <hr className="border-gray-600 my-6"/>

        {/* Bottom Section */}
        <div className="text-center text-gray-400 text-sm">
          <p>&copy; 2025 Project by <span className="text-white font-semibold">Team CODE CRAFTERS</span>.</p>
        </div>

      </div>
    </footer>
  );
}

export default Footer;
