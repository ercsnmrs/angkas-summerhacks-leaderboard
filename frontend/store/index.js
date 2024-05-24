import { v4 as uuidv4 } from "uuid";

export const state = () => ({
  microservice: {
    brokerURL: "https://api.banco.solutions",
    description: "Microservices is a good way to showcase expertise in  Application Programming Interface (API). APIs are the solution to serve as the glue that enables communication, abstraction, loose coupling, standardization, and security in microservices architecture, facilitating the development, integration, and maintenance of microservices-based applications.",
    features:{
      broker: {
        desc: "Request is sent to the broker service, base on the request it routes it to the proper service. In this case we just want to know if the broker service is active.",
      },
      mailer: {
        desc: "Request is sent to the broker service, redirects the email to our mailer service. Can check it on api.banco.solutions:8025",
      },
      rabbitLogger: {
        desc: "Request is sent to the broker service using RPC, from broker uses RabbitMQ to trigger an event. It will now try to process the log and will persist on MongoDB",
      },
      auth: {
        desc: "The mock user data is stored on Postgres. The request is sent to broker service then passes the request to the authentication service",
      },
      grpcLogger: {
        desc: "Unlike like the RabbitMQ Logger this logger uses GRPC instead of RPC.",
      },
    }

  },
  web3: {
    walletURL: "https://wallet.banco.solutions",
    chainURL: "https://chain.banco.solutions",
    description: "Using blockchain in web3 demos can help showcase its unique features and capabilities that align with the principles of decentralization, trust, transparency, data ownership and control, interoperability, and security, which are central to the web3 vision.",
    features:{
      scalability: {
        title: "Scalable",
        desc: "",
      },
      flexibility: {
        title: "Flexible",
        description: "",
      },
      maintainability: {
        title: "Maintainable",
        desc: "",
      },
    }

  },
  copyrightDate: new Date().getFullYear(),
  socialProfiles: [
    {
      id: uuidv4(),
      name: "GitHub",
      icon: "github",
      url: "https://github.com/ercsnmrs",
    },
    {
      id: uuidv4(),
      name: "Linkedin",
      icon: "linkedin",
      url: "https://linkedin.com/in/ercsnmrs",
    },
  ],

});

export const getters = {
  // @todo
};

export const mutations = {
  // @todo
};

export const actions = {
  // @todo
};
