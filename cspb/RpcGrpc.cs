// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: rpc.proto
#region Designer generated code

using System;
using System.Threading;
using System.Threading.Tasks;
using Grpc.Core;

namespace prpc {
  public static class RPCService
  {
    static readonly string __ServiceName = "prpc.RPCService";

    static readonly Marshaller<global::prpc.RPCRequest> __Marshaller_RPCRequest = Marshallers.Create((arg) => global::Google.Protobuf.MessageExtensions.ToByteArray(arg), global::prpc.RPCRequest.Parser.ParseFrom);
    static readonly Marshaller<global::prpc.RPCResponse> __Marshaller_RPCResponse = Marshallers.Create((arg) => global::Google.Protobuf.MessageExtensions.ToByteArray(arg), global::prpc.RPCResponse.Parser.ParseFrom);
    static readonly Marshaller<global::prpc.RPCStreamMessage> __Marshaller_RPCStreamMessage = Marshallers.Create((arg) => global::Google.Protobuf.MessageExtensions.ToByteArray(arg), global::prpc.RPCStreamMessage.Parser.ParseFrom);

    static readonly Method<global::prpc.RPCRequest, global::prpc.RPCResponse> __Method_RequestCall = new Method<global::prpc.RPCRequest, global::prpc.RPCResponse>(
        MethodType.Unary,
        __ServiceName,
        "RequestCall",
        __Marshaller_RPCRequest,
        __Marshaller_RPCResponse);

    static readonly Method<global::prpc.RPCStreamMessage, global::prpc.RPCStreamMessage> __Method_Streaming = new Method<global::prpc.RPCStreamMessage, global::prpc.RPCStreamMessage>(
        MethodType.DuplexStreaming,
        __ServiceName,
        "Streaming",
        __Marshaller_RPCStreamMessage,
        __Marshaller_RPCStreamMessage);

    // service descriptor
    public static global::Google.Protobuf.Reflection.ServiceDescriptor Descriptor
    {
      get { return global::prpc.Rpc.Descriptor.Services[0]; }
    }

    // client interface
    public interface IRPCServiceClient
    {
      global::prpc.RPCResponse RequestCall(global::prpc.RPCRequest request, Metadata headers = null, DateTime? deadline = null, CancellationToken cancellationToken = default(CancellationToken));
      global::prpc.RPCResponse RequestCall(global::prpc.RPCRequest request, CallOptions options);
      AsyncUnaryCall<global::prpc.RPCResponse> RequestCallAsync(global::prpc.RPCRequest request, Metadata headers = null, DateTime? deadline = null, CancellationToken cancellationToken = default(CancellationToken));
      AsyncUnaryCall<global::prpc.RPCResponse> RequestCallAsync(global::prpc.RPCRequest request, CallOptions options);
      AsyncDuplexStreamingCall<global::prpc.RPCStreamMessage, global::prpc.RPCStreamMessage> Streaming(Metadata headers = null, DateTime? deadline = null, CancellationToken cancellationToken = default(CancellationToken));
      AsyncDuplexStreamingCall<global::prpc.RPCStreamMessage, global::prpc.RPCStreamMessage> Streaming(CallOptions options);
    }

    // server-side interface
    public interface IRPCService
    {
      Task<global::prpc.RPCResponse> RequestCall(global::prpc.RPCRequest request, ServerCallContext context);
      Task Streaming(IAsyncStreamReader<global::prpc.RPCStreamMessage> requestStream, IServerStreamWriter<global::prpc.RPCStreamMessage> responseStream, ServerCallContext context);
    }

    // client stub
    public class RPCServiceClient : ClientBase, IRPCServiceClient
    {
      public RPCServiceClient(Channel channel) : base(channel)
      {
      }
      public global::prpc.RPCResponse RequestCall(global::prpc.RPCRequest request, Metadata headers = null, DateTime? deadline = null, CancellationToken cancellationToken = default(CancellationToken))
      {
        var call = CreateCall(__Method_RequestCall, new CallOptions(headers, deadline, cancellationToken));
        return Calls.BlockingUnaryCall(call, request);
      }
      public global::prpc.RPCResponse RequestCall(global::prpc.RPCRequest request, CallOptions options)
      {
        var call = CreateCall(__Method_RequestCall, options);
        return Calls.BlockingUnaryCall(call, request);
      }
      public AsyncUnaryCall<global::prpc.RPCResponse> RequestCallAsync(global::prpc.RPCRequest request, Metadata headers = null, DateTime? deadline = null, CancellationToken cancellationToken = default(CancellationToken))
      {
        var call = CreateCall(__Method_RequestCall, new CallOptions(headers, deadline, cancellationToken));
        return Calls.AsyncUnaryCall(call, request);
      }
      public AsyncUnaryCall<global::prpc.RPCResponse> RequestCallAsync(global::prpc.RPCRequest request, CallOptions options)
      {
        var call = CreateCall(__Method_RequestCall, options);
        return Calls.AsyncUnaryCall(call, request);
      }
      public AsyncDuplexStreamingCall<global::prpc.RPCStreamMessage, global::prpc.RPCStreamMessage> Streaming(Metadata headers = null, DateTime? deadline = null, CancellationToken cancellationToken = default(CancellationToken))
      {
        var call = CreateCall(__Method_Streaming, new CallOptions(headers, deadline, cancellationToken));
        return Calls.AsyncDuplexStreamingCall(call);
      }
      public AsyncDuplexStreamingCall<global::prpc.RPCStreamMessage, global::prpc.RPCStreamMessage> Streaming(CallOptions options)
      {
        var call = CreateCall(__Method_Streaming, options);
        return Calls.AsyncDuplexStreamingCall(call);
      }
    }

    // creates service definition that can be registered with a server
    public static ServerServiceDefinition BindService(IRPCService serviceImpl)
    {
      return ServerServiceDefinition.CreateBuilder(__ServiceName)
          .AddMethod(__Method_RequestCall, serviceImpl.RequestCall)
          .AddMethod(__Method_Streaming, serviceImpl.Streaming).Build();
    }

    // creates a new client
    public static RPCServiceClient NewClient(Channel channel)
    {
      return new RPCServiceClient(channel);
    }

  }
}
#endregion
