// See here for reference: https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-lambda-authorizer.html
module.exports.handler = async (event) => {
    const token = event.headers.authorization?.split("Bearer ")?.[1];
    const isAuthorized = token === process.env.AUTH_SECRET
    return { isAuthorized }
};
