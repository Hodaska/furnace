PROGRAM furnace_control
  VAR
    Start      AT %QX0.0 : BOOL;
    Alarm      AT %QX0.1 : BOOL;
    TempScaled AT %QW0   : INT;
    PowerOut   AT %QW1   : INT;
  END_VAR
  VAR
    Temperature : REAL;
  END_VAR

  Temperature := INT_TO_REAL(TempScaled) / 10.0;

  IF Temperature >= 100.0 THEN
    Alarm    := TRUE;
    Start    := FALSE;
    PowerOut := 0;
  ELSIF Alarm AND Temperature < 90.0 THEN
    Alarm    := FALSE;
    PowerOut := 0;
  ELSIF Start AND NOT Alarm THEN
    PowerOut := 100;
  ELSE
    PowerOut := 0;
  END_IF;

END_PROGRAM

CONFIGURATION Config0
  RESOURCE Res0 ON PLC
    TASK TaskMain(INTERVAL := T#100ms, PRIORITY := 0);
    PROGRAM Inst0 WITH TaskMain : furnace_control;
  END_RESOURCE
END_CONFIGURATION
