using Serilog.Events;
using System.Text;

namespace Learning.Infrastructure.Telemetry;

public static class LogHelper
{
    private static readonly HashSet<string> SystemLogCategories = new(StringComparer.OrdinalIgnoreCase)
    {
        "Microsoft",
        "System"
    };

    public static readonly HashSet<string> IgnoreLogProperties = new(StringComparer.OrdinalIgnoreCase)
    {
        "MessageTemplate",
        "RequestPath",
        "RequestId"
    };

    public static readonly Dictionary<LogEventLevel, string> LogLevelMap = new()
    {
        [LogEventLevel.Debug] = "DEBUG",
        [LogEventLevel.Information] = "INFO",
        [LogEventLevel.Warning] = "WARN",
        [LogEventLevel.Error] = "ERROR",
        [LogEventLevel.Fatal] = "FATAL"
    };

    public static bool IsSystemLogCategory(string categoryName)
    {
        return SystemLogCategories.Any(categoryName.Contains);
    }

    public static string ToSnakeCase(string input)
    {
        if (string.IsNullOrWhiteSpace(input)) return input;

        var sb = new StringBuilder();
        for (int i = 0; i < input.Length; i++)
        {
            char c = input[i];
            if (char.IsUpper(c) && i > 0)
                sb.Append('_');
            sb.Append(char.ToLowerInvariant(c));
        }
        return sb.ToString();
    }
}